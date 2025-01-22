package proxy

import (
	"runtime/debug"
)

// Debug is an interface that provides a proxy of the methods of debug.
type Debug interface {
	ReadBuildInfo() (info *debug.BuildInfo, ok bool)
}

// debugProxy is a proxy struct that implements the Debug interface.
type debugProxy struct{}

// NewDebug returns a new instance of the Debug interface.
func NewDebug() Debug {
	return &debugProxy{}
}

// ReadBuildInfo is a proxy method that calls the ReadBuildInfo method of the debug.
func (*debugProxy) ReadBuildInfo() (info *debug.BuildInfo, ok bool) {
	return debug.ReadBuildInfo()
}
