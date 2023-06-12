package expr

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"gonum.org/v1/gonum/graph/simple"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/expr/mathexp"
	"github.com/grafana/grafana/pkg/expr/ml"
	"github.com/grafana/grafana/pkg/plugins/httpresponsesender"
	"github.com/grafana/grafana/pkg/services/pluginsintegration/plugincontext"
)

var errMLPluginDoesNotExist = errors.New("expression type Machine Learning is not supported. Plugin 'grafana-ml-app' must be enabled")

type MLNode struct {
	baseNode
	command   ml.Command
	TimeRange TimeRange
	request   *Request
}

// NodeType returns the data pipeline node type.
func (m *MLNode) NodeType() NodeType {
	return TypeMlNode
}

// Execute runs ml.Command command and converts the response to mathexp.Results
func (m *MLNode) Execute(ctx context.Context, now time.Time, _ mathexp.Vars, s *Service) (r mathexp.Results, e error) {
	logger := logger.FromContext(ctx).New("datasourceType", TypeMlNode.String(), "queryRefId", m.refID)
	var result mathexp.Results
	timeRange := m.TimeRange.AbsoluteTime(now)

	pCtx, err := s.pCtxProvider.Get(ctx, ml.PluginID, m.request.User, m.request.OrgId)
	if err != nil {
		if errors.Is(err, plugincontext.ErrPluginNotFound) {
			return result, errMLPluginDoesNotExist
		}
		return result, fmt.Errorf("failed to get plugin settings: %w", err)
	}

	responseType := "unknown"
	respStatus := "success"
	defer func() {
		if e != nil {
			responseType = "error"
			respStatus = "failure"
		}
		logger.Debug("Data source queried", "responseType", responseType)
		useDataplane := strings.HasPrefix("dataplane-", responseType)
		s.metrics.dsRequests.WithLabelValues(respStatus, fmt.Sprintf("%t", useDataplane)).Inc()
	}()

	resp, err := m.command.Execute(ctx, timeRange.From, timeRange.To, func(path string, payload []byte) ([]byte, error) {
		crReq := &backend.CallResourceRequest{
			PluginContext: pCtx,
			Path:          path,
			Method:        http.MethodPost,
			URL:           path,
			Headers:       make(map[string][]string, len(m.request.Headers)),
			Body:          payload,
		}

		for key, val := range m.request.Headers {
			crReq.Headers[key] = []string{val}
		}

		resp := response.CreateNormalResponse(make(http.Header), nil, 0)
		httpSender := httpresponsesender.New(resp)
		err = s.pluginsClient.CallResource(ctx, crReq, httpSender)
		if err != nil {
			return nil, err
		}

		if resp.Status() >= 200 && resp.Status() < 300 {
			return resp.Body(), nil
		}
		return nil, fmt.Errorf("failed to send a POST request to plugin %s via proxy by path %s, status code: %v, msg:%s", ml.PluginID, path, resp.Status(), resp.Body())
	})

	if err != nil {
		return result, QueryError{
			RefID: m.refID,
			Err:   err,
		}
	}

	responseType, result, err = queryDataResponseToResults(ctx, resp, m.refID, ml.PluginID, s)
	return result, err
}

func (s *Service) buildMlNode(dp *simple.DirectedGraph, rn *rawNode, req *Request) (Node, error) {
	// make sure that unattached context has bounds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, exist := s.plugins.Plugin(ctx, ml.PluginID); !exist {
		return nil, errMLPluginDoesNotExist
	}
	if rn.TimeRange == nil {
		return nil, errors.New("time range must be specified")
	}

	cmd, err := ml.UnmarshalCommand(rn.Query, s.cfg.AppURL)
	if err != nil {
		return nil, err
	}

	return &MLNode{
		baseNode: baseNode{
			id:    dp.NewNode().ID(),
			refID: rn.RefID,
		},
		TimeRange: rn.TimeRange,
		command:   cmd,
		request:   req,
	}, nil
}
