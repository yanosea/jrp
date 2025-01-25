package proxy

import (
	"encoding/json"
)

// Json is an interface that provides a proxy of the methods of encoding/json.
type Json interface {
	Marshal(v interface{}) ([]byte, error)
}

// jsonProxy is a proxy struct that implements the Json interface.
type jsonProxy struct{}

// NewJson returns a new instance of the Json interface.
func NewJson() Json {
	return &jsonProxy{}
}

// Marshal returns the JSON encoding of v.
func (j *jsonProxy) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
