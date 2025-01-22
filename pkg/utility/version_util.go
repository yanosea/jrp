package utility

import (
	"github.com/yanosea/jrp/pkg/proxy"
)

// VersionUtil is an interface that provides the version of the application.
type VersionUtil interface {
	GetVersion(version string) string
}

// versionUtil is a struct that implements the VersionUtil interface.
type versionUtil struct {
	debug proxy.Debug
}

// NewVersionUtil returns a new instance of the versionUtil struct.
func NewVersionUtil(
	debug proxy.Debug,
) VersionUtil {
	return &versionUtil{
		debug: debug,
	}
}

// GetVersion returns the version of the application.
func (v *versionUtil) GetVersion(version string) string {
	// if version is embedded, return it
	if version != "" {
		return version
	}

	if i, ok := v.debug.ReadBuildInfo(); !ok {
		return "unknown"
	} else if i.Main.Version != "" {
		return i.Main.Version
	} else {
		return "dev"
	}
}
