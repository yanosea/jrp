package buildinfo

import (
	"runtime/debug"
)

type BuildInfoProvider interface {
	ReadBuildInfo() (*debug.BuildInfo, bool)
}

type RealBuildInfoProvider struct{}

func (RealBuildInfoProvider) ReadBuildInfo() (*debug.BuildInfo, bool) {
	return debug.ReadBuildInfo()
}
