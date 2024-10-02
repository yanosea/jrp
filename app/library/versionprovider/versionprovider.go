package versionprovider

import (
	"github.com/yanosea/jrp/app/proxy/debug"
)

// VersionProvidable is an interface for VersionProvider.
type VersionProvidable interface {
	GetVersion(emmbeddedVersion string) string
}

// VersionProvider is a struct that implements VersionProvidable.
type VersionProvider struct {
	DebugProxy debugproxy.Debug
}

// New is a constructor of VersionProvider.
func New(debugProxy debugproxy.Debug) *VersionProvider {
	return &VersionProvider{
		DebugProxy: debugProxy,
	}
}

// GetVersion gets the version from build info if version is not embedded.
func (v *VersionProvider) GetVersion(embeddedVersion string) string {
	// if version is embedded, return it
	if embeddedVersion != "" {
		return embeddedVersion
	}

	i, ok := v.DebugProxy.ReadBuildInfo()
	if !ok {
		// if reading build info fails, return unknown
		return "unknown"
	}
	if i.FieldBuildInfo.Main.Version == "" || i.FieldBuildInfo.Main.Version == "(devel)" {
		// if version from build info is empty, return dev
		return "devel"
	}

	// return version from build info
	return i.FieldBuildInfo.Main.Version
}
