package logic

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/util"
)

func Download(u UserProvider, hc HTTPClient, fs FileSystem) error {
	dbFileDirPath, err := u.GetDBFileDirPath()
	if err != nil {
		return err
	}

	if _, err := fs.Stat(dbFileDirPath); os.IsNotExist(err) {
		fs.MkdirAll(dbFileDirPath, 0755)
	}

	dbFilePath := filepath.Join(dbFileDirPath, constant.WNJPN_DB_FILE_NAME)
	if _, err := fs.Stat(dbFilePath); os.IsNotExist(err) {
		resp, err := hc.Get(constant.WNJPN_DB_ARCHIVE_FILE_URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		tempFilePath := filepath.Join(fs.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
		out, err := fs.Create(tempFilePath)
		if err != nil {
			return err
		}
		defer out.Close()
		if _, err := io.Copy(out, resp.Body); err != nil {
			return err
		}
		if _, err := out.Seek(0, 0); err != nil {
			return err
		}

		gz, err := gzip.NewReader(out)
		if err != nil {
			return err
		}
		defer gz.Close()

		f, err := fs.Create(dbFilePath)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, gz); err != nil {
			return err
		}

		if err := fs.Remove(tempFilePath); err != nil {
			return err
		}
		fmt.Println(constant.DOWNLOAD_MESSAGE_SUCCEEDED)
	} else {
		fmt.Println(constant.DOWNLOAD_MESSAGE_ALREADY_DOWNLOADED)
	}

	return nil
}

// for testing
var userProviderInstance = &DefaultUserProvider{}
var httpClientInstance = &DefaultHTTPClient{}
var fsInstance = &DefaultFileSystem{}

type UserProvider interface {
	GetDBFileDirPath() (string, error)
}

type DefaultUserProvider struct{}

func (p DefaultUserProvider) GetDBFileDirPath() (string, error) {
	user := &util.DefaultUserProvider{}
	return util.GetDBFileDirPath(user)
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type DefaultHTTPClient struct{}

func (c *DefaultHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	TempDir() string
	Create(name string) (*os.File, error)
	Remove(name string) error
}

type DefaultFileSystem struct{}

func (fs DefaultFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs DefaultFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs DefaultFileSystem) TempDir() string {
	return os.TempDir()
}

func (fs DefaultFileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (fs DefaultFileSystem) Remove(name string) error {
	return os.Remove(name)
}
