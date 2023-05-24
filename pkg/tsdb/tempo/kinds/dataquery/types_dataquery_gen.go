// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     PluginGoTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package dataquery

// Defines values for SearchStreamingState.
const (
	SearchStreamingStateDone      SearchStreamingState = "done"
	SearchStreamingStateError     SearchStreamingState = "error"
	SearchStreamingStatePending   SearchStreamingState = "pending"
	SearchStreamingStateStreaming SearchStreamingState = "streaming"
)

// Defines values for TempoQueryFiltersScope.
const (
	TempoQueryFiltersScopeResource TempoQueryFiltersScope = "resource"
	TempoQueryFiltersScopeSpan     TempoQueryFiltersScope = "span"
	TempoQueryFiltersScopeUnscoped TempoQueryFiltersScope = "unscoped"
)

// Defines values for TempoQueryType.
const (
	TempoQueryTypeClear         TempoQueryType = "clear"
	TempoQueryTypeNativeSearch  TempoQueryType = "nativeSearch"
	TempoQueryTypeSearch        TempoQueryType = "search"
	TempoQueryTypeServiceMap    TempoQueryType = "serviceMap"
	TempoQueryTypeTraceId       TempoQueryType = "traceId"
	TempoQueryTypeTraceql       TempoQueryType = "traceql"
	TempoQueryTypeTraceqlSearch TempoQueryType = "traceqlSearch"
	TempoQueryTypeUpload        TempoQueryType = "upload"
)

// Defines values for TraceqlFilterScope.
const (
	TraceqlFilterScopeResource TraceqlFilterScope = "resource"
	TraceqlFilterScopeSpan     TraceqlFilterScope = "span"
	TraceqlFilterScopeUnscoped TraceqlFilterScope = "unscoped"
)

// Defines values for TraceqlSearchScope.
const (
	TraceqlSearchScopeResource TraceqlSearchScope = "resource"
	TraceqlSearchScopeSpan     TraceqlSearchScope = "span"
	TraceqlSearchScopeUnscoped TraceqlSearchScope = "unscoped"
)

// The state of the TraceQL streaming search query
type SearchStreamingState string

// TempoDataQuery defines model for TempoDataQuery.
type TempoDataQuery = map[string]interface{}

// TempoQuery defines model for TempoQuery.
type TempoQuery struct {
	// For mixed data sources the selected datasource is on the query level.
	// For non mixed scenarios this is undefined.
	// TODO find a better way to do this ^ that's friendly to schema
	// TODO this shouldn't be unknown but DataSourceRef | null
	Datasource *interface{} `json:"datasource,omitempty"`
	Filters    []struct {
		// Uniquely identify the filter, will not be used in the query generation
		Id string `json:"id"`

		// The operator that connects the tag to the value, for example: =, >, !=, =~
		Operator *string `json:"operator,omitempty"`

		// The scope of the filter, can either be unscoped/all scopes, resource or span
		Scope *TempoQueryFiltersScope `json:"scope,omitempty"`

		// The tag for the search filter, for example: .http.status_code, .service.name, status
		Tag *string `json:"tag,omitempty"`

		// The value for the search filter
		Value *interface{} `json:"value,omitempty"`

		// The type of the value, used for example to check whether we need to wrap the value in quotes when generating the query
		ValueType *string `json:"valueType,omitempty"`
	} `json:"filters"`

	// Hide true if query is disabled (ie should not be returned to the dashboard)
	// Note this does not always imply that the query should not be executed since
	// the results from a hidden query may be used as the input to other queries (SSE etc)
	Hide *bool `json:"hide,omitempty"`

	// Defines the maximum number of traces that are returned from Tempo
	Limit *int64 `json:"limit,omitempty"`

	// Define the maximum duration to select traces. Use duration format, for example: 1.2s, 100ms
	MaxDuration *string `json:"maxDuration,omitempty"`

	// Define the minimum duration to select traces. Use duration format, for example: 1.2s, 100ms
	MinDuration *string `json:"minDuration,omitempty"`

	// TraceQL query or trace ID
	Query string `json:"query"`

	// Specify the query flavor
	// TODO make this required and give it a default
	QueryType *string `json:"queryType,omitempty"`

	// A unique identifier for the query within the list of targets.
	// In server side expressions, the refId is used as a variable name to identify results.
	// By default, the UI will assign A->Z; however setting meaningful names may be useful.
	RefId string `json:"refId"`

	// Logfmt query to filter traces by their tags. Example: http.status_code=200 error=true
	Search *string `json:"search,omitempty"`

	// Filters to be included in a PromQL query to select data for the service graph. Example: {client="app",service="app"}
	ServiceMapQuery *string `json:"serviceMapQuery,omitempty"`

	// Query traces by service name
	ServiceName *string `json:"serviceName,omitempty"`

	// Query traces by span name
	SpanName *string `json:"spanName,omitempty"`

	// Use the streaming API to get partial results as they are available
	Streaming *bool `json:"streaming,omitempty"`
}

// The scope of the filter, can either be unscoped/all scopes, resource or span
type TempoQueryFiltersScope string

// TempoQueryType search = Loki search, nativeSearch = Tempo search for backwards compatibility
type TempoQueryType string

// TraceqlFilter defines model for TraceqlFilter.
type TraceqlFilter struct {
	// Uniquely identify the filter, will not be used in the query generation
	Id string `json:"id"`

	// The operator that connects the tag to the value, for example: =, >, !=, =~
	Operator *string `json:"operator,omitempty"`

	// The scope of the filter, can either be unscoped/all scopes, resource or span
	Scope *TraceqlFilterScope `json:"scope,omitempty"`

	// The tag for the search filter, for example: .http.status_code, .service.name, status
	Tag *string `json:"tag,omitempty"`

	// The value for the search filter
	Value *interface{} `json:"value,omitempty"`

	// The type of the value, used for example to check whether we need to wrap the value in quotes when generating the query
	ValueType *string `json:"valueType,omitempty"`
}

// The scope of the filter, can either be unscoped/all scopes, resource or span
type TraceqlFilterScope string

// TraceqlSearchScope static fields are pre-set in the UI, dynamic fields are added by the user
type TraceqlSearchScope string
