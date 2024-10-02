package httpproxy

import (
	"net/http"
)

// Http is an interface for http.
type Http interface {
	Get(url string) (*ResponseInstance, error)
}

// HttpProxy is a struct that implements Http.
type HttpProxy struct{}

// New is a constructor for HttpProxy.
func New() Http {
	return &HttpProxy{}
}

// Get is a proxy for http.Get.
func (*HttpProxy) Get(url string) (*ResponseInstance, error) {
	resp, err := http.Get(url)
	return &ResponseInstance{FieldResponse: resp}, err
}
