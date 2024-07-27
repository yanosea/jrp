package logic

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
)

type MockUserProvider struct {
	DBFileDirPath       string
	GetDBFileDirPathErr error
}

func (p MockUserProvider) GetDBFileDirPath() (string, error) {
	if p.GetDBFileDirPathErr != nil {
		return "", p.GetDBFileDirPathErr
	}
	return p.DBFileDirPath, nil
}

type MockHTTPClient struct {
	Response *http.Response
	Error    error
}

func (c MockHTTPClient) Get(url string) (*http.Response, error) {
	return c.Response, c.Error
}

type MockFileSystem struct {
	StatErr     error
	MkdirAllErr error
	CreateFile  *os.File
	CreateErr   error
	TempDirPath string
	RemoveErr   error
}

func (fs MockFileSystem) Stat(name string) (os.FileInfo, error) {
	return nil, fs.StatErr
}

func (fs MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return fs.MkdirAllErr
}

func (fs MockFileSystem) TempDir() string {
	return fs.TempDirPath
}

func (fs MockFileSystem) Create(name string) (*os.File, error) {
	return fs.CreateFile, fs.CreateErr
}

func (fs MockFileSystem) Remove(name string) error {
	return fs.RemoveErr
}

func TestDownload(t *testing.T) {
	const testDBFileDirPath = "/mock/path/to/db"
	const testTempFilePath = "/mock/temp/dir/mockedfile.gz"
	const testDBFilePath = "/mock/path/to/dbfile"

	tests := []struct {
		name         string
		userProvider UserProvider
		httpClient   HTTPClient
		fs           FileSystem
		wantError    bool
	}{
		{
			name:         "Success",
			userProvider: MockUserProvider{DBFileDirPath: testDBFileDirPath},
			httpClient: MockHTTPClient{
				Response: &http.Response{
					Body: io.NopCloser(bytes.NewBuffer([]byte("gzipdata"))),
				},
			},
			fs: MockFileSystem{
				StatErr:     os.ErrNotExist,
				TempDirPath: "/mock/temp/dir",
				CreateFile:  createMockGzipFile("mock gzipped data"),
			},
			wantError: false,
		},
		{
			name:         "UserProviderError",
			userProvider: MockUserProvider{GetDBFileDirPathErr: errors.New("GetDBFileDirPath error")},
			httpClient:   MockHTTPClient{},
			fs:           MockFileSystem{},
			wantError:    true,
		},
		{
			name:         "FileExists",
			userProvider: MockUserProvider{DBFileDirPath: testDBFileDirPath},
			httpClient:   MockHTTPClient{},
			fs: MockFileSystem{
				StatErr: nil,
			},
			wantError: false,
		},
		{
			name:         "MkdirAllError",
			userProvider: MockUserProvider{DBFileDirPath: testDBFileDirPath},
			httpClient:   MockHTTPClient{},
			fs: MockFileSystem{
				StatErr:     os.ErrNotExist,
				MkdirAllErr: errors.New("MkdirAll error"),
			},
			wantError: true,
		},
		{
			name:         "HttpClientError",
			userProvider: MockUserProvider{DBFileDirPath: testDBFileDirPath},
			httpClient: MockHTTPClient{
				Error: errors.New("HTTP client error"),
			},
			fs: MockFileSystem{
				StatErr: os.ErrNotExist,
			},
			wantError: true,
		},
		{
			name:         "FileCreateError",
			userProvider: MockUserProvider{DBFileDirPath: testDBFileDirPath},
			httpClient: MockHTTPClient{
				Response: &http.Response{
					Body: io.NopCloser(bytes.NewBuffer([]byte("gzipdata"))),
				},
			},
			fs: MockFileSystem{
				StatErr:     os.ErrNotExist,
				TempDirPath: "/mock/temp/dir",
				CreateErr:   errors.New("Create error"),
			},
			wantError: true,
		},
		{
			name:         "GzipReaderError",
			userProvider: MockUserProvider{DBFileDirPath: testDBFileDirPath},
			httpClient: MockHTTPClient{
				Response: &http.Response{
					Body: io.NopCloser(bytes.NewBuffer([]byte("invalid gzip data"))),
				},
			},
			fs: MockFileSystem{
				StatErr:     os.ErrNotExist,
				TempDirPath: "/mock/temp/dir",
				CreateFile:  createMockGzipFile("invalid gzip data"),
			},
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Download(tc.userProvider, tc.httpClient, tc.fs)
			if (err != nil) != tc.wantError {
				t.Errorf("Download() error = %v, wantError %v", err, tc.wantError)
			}
		})
	}
}

func createMockGzipFile(content string) *os.File {
	data := bytes.NewBuffer([]byte{})

	writer := gzip.NewWriter(data)
	writer.Write([]byte(content))
	writer.Close()

	tmpFile, _ := os.CreateTemp("", "test*.gz")
	tmpFile.Write(data.Bytes())
	tmpFile.Seek(0, 0)
	return tmpFile
}
