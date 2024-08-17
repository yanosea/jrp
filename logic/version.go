package logic

import (
	"github.com/yanosea/jrp/internal/buildinfo"
)

type VersionGetter interface {
	GetVersion(string) string
}

type JrpVersionGetter struct {
	BuildInfoProvider buildinfo.BuildInfoProvider
}

func NewJrpVersionGetter(b buildinfo.BuildInfoProvider) *JrpVersionGetter {
	return &JrpVersionGetter{
		BuildInfoProvider: b,
	}
}

func (b *JrpVersionGetter) GetVersion(v string) string {
	// if version is embedded, return it
	if v != "" {
		return v
	}

	if i, ok := b.BuildInfoProvider.ReadBuildInfo(); !ok {
		// return unknown
		return "unknown"
	} else {
		// return version from build info
		return i.Main.Version
	}
}
