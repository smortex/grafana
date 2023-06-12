package ml

import (
	"fmt"
)

// DatasourceType is the string constant used as the datasource when the property is in Datasource.Type.
// Type in requests is used to identify what type of data source plugin the request belongs to.
const DatasourceType = "__ml__"

// DatasourceUID is the string constant used as the datasource name in requests
// to identify it as an expression command when use in Datasource.UID.
const DatasourceUID = DatasourceType

// IsDataSource checks if the uid points to ML node query
func IsDataSource(uid string) bool {
	return uid == DatasourceUID
}

func readValue[T any](query map[string]interface{}, key string) (T, error) {
	var result T
	v, ok := query[key]
	if !ok {
		return result, fmt.Errorf("required field '%s' is missing", key)
	}
	result, ok = v.(T)
	if !ok {
		return result, fmt.Errorf("field '%s' has type %T but expected string", key, v)
	}
	return result, nil
}

func readOptionalValue[T any](query map[string]interface{}, key string) (*T, error) {
	var result T
	v, ok := query[key]
	if !ok {
		return nil, fmt.Errorf("required field '%s' is missing", key)
	}
	result, ok = v.(T)
	if !ok {
		return nil, fmt.Errorf("field '%s' has type %T but expected string", key, v)
	}
	return &result, nil
}
