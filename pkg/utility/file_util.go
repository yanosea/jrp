package utility

import (
	"io"
	"path/filepath"

	"github.com/yanosea/jrp/pkg/proxy"
)

// FileUtil is an interface that contains the utility functions for file operations.
type FileUtil interface {
	ExtractGzFile(gzFilePath, destDir string) error
	GetXDGDataHome() (string, error)
	HideFile(filePath string) (string, error)
	IsExist(name string) bool
	MkdirIfNotExist(dirPath string) error
	RemoveAll(path string) error
	SaveToTempFile(body io.Reader, fileName string) (string, error)
	UnhideFile(filePath string) error
}

// fileUtil is a struct that contains the utility functions for file operations.
type fileUtil struct {
	gzip proxy.Gzip
	io   proxy.Io
	os   proxy.Os
}

// NewFileUtil returns a new instance of the FileUtil struct.
func NewFileUtil(
	gzip proxy.Gzip,
	io proxy.Io,
	os proxy.Os,
) FileUtil {
	return &fileUtil{
		gzip: gzip,
		io:   io,
		os:   os,
	}
}

// ExtractGzFile extracts a gzipped file to the destination directory.
func (f *fileUtil) ExtractGzFile(gzFilePath, destFilePath string) error {
	var deferErr error
	gzFile, err := f.os.Open(gzFilePath)
	if err != nil {
		return err
	}
	defer func() {
		deferErr = gzFile.Close()
	}()

	gzReader, err := f.gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer func() {
		deferErr = gzReader.Close()
	}()

	destFile, err := f.os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer func() {
		deferErr = destFile.Close()
	}()

	if _, err := f.io.Copy(destFile, gzReader); err != nil {
		return err
	}

	return deferErr
}

// GetXDGDataHome returns the XDG data home directory.
func (f *fileUtil) GetXDGDataHome() (string, error) {
	xdgDataHome := f.os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		homeDir, err := f.os.UserHomeDir()
		if err != nil {
			return "", err
		}

		xdgDataHome = filepath.Join(homeDir, ".local", "share")
	}

	return xdgDataHome, nil
}

// HideFile hides the file by adding a dot prefix to the file name.
func (f *fileUtil) HideFile(filePath string) (string, error) {
	hiddenFilePath := filepath.Join(filepath.Dir(filePath), "."+filepath.Base(filePath))
	if err := f.os.Rename(filePath, hiddenFilePath); err != nil {
		return "", err
	}

	return hiddenFilePath, nil
}

// IsExist checks if the file or directory exists.
func (f *fileUtil) IsExist(name string) bool {
	_, err := f.os.Stat(name)
	return !f.os.IsNotExist(err)
}

// MkdirIfNotExist creates a directory if it does not exist.
func (f *fileUtil) MkdirIfNotExist(dirPath string) error {
	if _, err := f.os.Stat(dirPath); f.os.IsNotExist(err) {
		if err := f.os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

// RemoveAll removes path and any children it contains.
func (f *fileUtil) RemoveAll(path string) error {
	return f.os.RemoveAll(path)
}

// SaveToTempFile saves the body to a temporary file.
func (f *fileUtil) SaveToTempFile(body io.Reader, fileName string) (string, error) {
	var deferErr error
	tempFilePath := filepath.Join(f.os.TempDir(), fileName)

	file, err := f.os.Create(tempFilePath)
	if err != nil {
		return "", err
	}
	defer func() {
		deferErr = file.Close()
	}()

	if _, err := f.io.Copy(file, body); err != nil {
		return "", err
	}

	return tempFilePath, deferErr
}

// UnhideFile unhides the file by removing the dot prefix from the file name.
func (f *fileUtil) UnhideFile(hiddenFilePath string) error {
	filePath := filepath.Join(filepath.Dir(hiddenFilePath), filepath.Base(hiddenFilePath)[1:])
	if err := f.os.Rename(hiddenFilePath, filePath); err != nil {
		return err
	}

	return nil
}
