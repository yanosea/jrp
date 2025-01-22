package proxy

import (
	"io"
	"net/http"
)

// Http is an interface that provides a proxy of the methods of http.
type Http interface {
	Get(url string) (Response, error)
}

// httpProxy is a proxy struct that implements the Http interface.
type httpProxy struct{}

// NewHttp returns a new instance of the Http interface.
func NewHttp() Http {
	return &httpProxy{}
}

// Get issues a GET to the specified URL.
func (httpProxy) Get(url string) (Response, error) {
	response, err := http.Get(url)
	return &responseProxy{response: response}, err
}

// Response is an interface that provides a proxy of the methods of http.Response.
type Response interface {
	Close() error
	GetBody() ReadCloser
}

// responseProxy is a proxy struct that implements the Response interface.
type responseProxy struct {
	response *http.Response
}

// Close closes the response body.
func (r *responseProxy) Close() error {
	return r.response.Body.Close()
}

// GetBody returns the response body.
func (r *responseProxy) GetBody() ReadCloser {
	return &readCloserProxy{readCloser: r.response.Body}
}

// ReadCloser is an interface that provides a proxy of the methods of io.ReadCloser.
type ReadCloser interface {
	Close() error
	Read(p []byte) (n int, err error)
}

// readCloserProxy is a proxy struct that implements the ReadCloser interface.
type readCloserProxy struct {
	readCloser io.ReadCloser
}

// Close closes the response body.
func (r *readCloserProxy) Close() error {
	return r.readCloser.Close()
}

// Read reads data from the response body.
func (r *readCloserProxy) Read(p []byte) (n int, err error) {
	return r.readCloser.Read(p)
}
