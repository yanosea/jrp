package debugproxy

import (
	"runtime/debug"
)

// BuildInfoInstanceInterface is  an interface for debug.BuildInfo.
type BuildInfoInstanceInterface interface{}

// BuildInfoInstance is a struct that implements BuildInfoInstanceInterface.
type BuildInfoInstance struct {
	FieldBuildInfo *debug.BuildInfo
}
