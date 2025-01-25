package utility

import (
	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// JsonUtil is an interface that contains the utility functions for JSON.
type JsonUtil interface {
	Marshal(v interface{}) ([]byte, error)
}

// jsonUtil is a struct that contains the utility functions for JSON.
type jsonUtil struct {
	json proxy.Json
}

// NewJsonUtil returns a new instance of the JsonUtil struct.
func NewJsonUtil(json proxy.Json) JsonUtil {
	return &jsonUtil{
		json: json,
	}
}

// Marshal marshals v into JSON.
func (ju *jsonUtil) Marshal(v interface{}) ([]byte, error) {
	return ju.json.Marshal(v)
}
