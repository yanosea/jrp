package testenvsetter

import (
	"github.com/yanosea/jrp/app/library/dbfiledirpathprovider"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/os"
)

// TestEnvSetterInterface is an interface for setting test environment.
type TestEnvSetterInterface interface {
	SetTestEnv() error
	UnsetTestEnv() error
}

// TestEnvSetter is a struct for setting test environment.
type TestEnvSetter struct {
	FilePathProxy filepathproxy.FilePath
	OsProxy       osproxy.Os
}

// New is a constructor for TestEnvSetter.
func New(filePathProxy filepathproxy.FilePath, osProxy osproxy.Os) *TestEnvSetter {
	return &TestEnvSetter{
		FilePathProxy: filePathProxy,
		OsProxy:       osProxy,
	}
}

// SetTestEnv sets test environment.
func (t *TestEnvSetter) SetTestEnv() error {
	//get temporary directory
	tempDir := t.OsProxy.TempDir()
	jrpTempDir := t.FilePathProxy.Join(tempDir, "jrp")
	//set JRP_ENV_WNJPN_DB_FILE_DIR to temporary directory
	if err := t.OsProxy.Setenv(dbfiledirpathprovider.JRP_ENV_WNJPN_DB_FILE_DIR, jrpTempDir); err != nil {
		return err
	}
	//set JRP_ENV_JRP_DB_FILE_DIR to temporary directory
	if err := t.OsProxy.Setenv(dbfiledirpathprovider.JRP_ENV_JRP_DB_FILE_DIR, jrpTempDir); err != nil {
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
