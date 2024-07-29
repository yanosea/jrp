package logic

import (
	"net/http"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}
type DefaultHttpClient struct{}

func (DefaultHttpClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}
