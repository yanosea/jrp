package utility

import (
	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// DownloadUtil is an interface that contains the utility functions for downloading files.
type DownloadUtil interface {
	Download(url string) (proxy.Response, error)
}

// downloadUtil is a struct that contains the utility functions for downloading files.
type downloadUtil struct {
	http proxy.Http
}

// NewDownloadUtil returns a new instance of the DownloadUtil struct.
func NewDownloadUtil(http proxy.Http) DownloadUtil {
	return &downloadUtil{
		http: http,
	}
}

// Download downloads a file from the given URL.
func (d *downloadUtil) Download(url string) (proxy.Response, error) {
	if res, err := d.http.Get(url); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}
