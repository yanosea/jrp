package testutility

import (
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/proxy/os"
)

// TestEnvSetterInterface is an interface for setting test environment.
type TestEnvSetterInterface interface {
}

// TestEnvSetter is a struct for setting test environment.
type TestEnvSetter struct {
	OsProxy osproxy.Os
}

// NewTestEnvSetter is a constructor for TestEnvSetter.
func NewTestEnvSetter(osProxy osproxy.Os) *TestEnvSetter {
	return &TestEnvSetter{
		OsProxy: osProxy,
	}
}

// SetTestEnv sets test environment.
func (t *TestEnvSetter) SetTestEnv() error {
	//get temporary directory
	tempDir := t.OsProxy.TempDir()
	//set JRP_ENV_WNJPN_DB_FILE_DIR to temporary directory
	if err := t.OsProxy.Setenv(dbfiledirpathprovider.JRP_ENV_WNJPN_DB_FILE_DIR, tempDir); err != nil {
		return err
	}
	//set JRP_ENV_JRP_DB_FILE_DIR to temporary directory
	if err := t.OsProxy.Setenv(dbfiledirpathprovider.JRP_ENV_JRP_DB_FILE_DIR, tempDir); err != nil {
		return err
	}

	return nil
}

// UnsetTestEnv unsets test environment.
func (t *TestEnvSetter) UnsetTestEnv() error {
	//unset JRP_ENV_WNJPN_DB_FILE_DIR
	if err := t.OsProxy.Unsetenv(dbfiledirpathprovider.JRP_ENV_WNJPN_DB_FILE_DIR); err != nil {
		return err
	}
	//unset JRP_ENV_JRP_DB_FILE_DIR
	if err := t.OsProxy.Unsetenv(dbfiledirpathprovider.JRP_ENV_JRP_DB_FILE_DIR); err != nil {
		return err
	}

	return nil
}
