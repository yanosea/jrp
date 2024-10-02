package httpproxy

import (
	"net/http"
)

// ResponseInstanceInterface is an interface for http.Response.
type ResponseInstanceInterface interface{}

// ResponseInstance is a struct that implements ResponseInstanceInterface.
type ResponseInstance struct {
	FieldResponse *http.Response
}
