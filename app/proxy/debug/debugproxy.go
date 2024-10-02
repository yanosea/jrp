package debugproxy

import (
	"runtime/debug"
)

// Debug is an interface for debug.
type Debug interface {
	ReadBuildInfo() (*BuildInfoInstance, bool)
}

// DebugProxy is a struct that implements Debug.
type DebugProxy struct{}

// New is a constructor for DebugProxy.
func New() Debug {
	return &DebugProxy{}
}

// ReadBuildInfo is a proxy for debug.ReadBuildInfo.
func (*DebugProxy) ReadBuildInfo() (*BuildInfoInstance, bool) {
	buildInfo, ok := debug.ReadBuildInfo()
	return &BuildInfoInstance{FieldBuildInfo: buildInfo}, ok
}
