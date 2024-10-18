package testutility

import (
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/strings"
)

// FileHiderInterface is an interface for hiding and restoring file.
type FileHiderInterface interface {
	HideFile(filePath string) (int, error)
	RestoreFile(index int) error
}

// FileHider is a struct for hiding and restoring file.
type FileHider struct {
	FilePathProxy filepathproxy.FilePath
	OsProxy       osproxy.Os
	StringsProxy  stringsproxy.Strings
	HiddenFiles   []string
}

// NewFileMover is a constructor for FileMover.
func NewFileMover(filePathProxy filepathproxy.FilePath, osProxy osproxy.Os, stringsProxy stringsproxy.Strings) *FileHider {
	return &FileHider{
		FilePathProxy: filePathProxy,
		OsProxy:       osProxy,
		StringsProxy:  stringsProxy,
		HiddenFiles:   []string{},
	}
}

// HideFile hides file by renaming it to hidden file and returns hidden file's index.
func (f *FileHider) HideFile(filePath string) (int, error) {
	// check if file exists
	if _, err := f.OsProxy.Stat(filePath); f.OsProxy.IsNotExist(err) {
		return -1, err
	}
	// extract the directory and file name
	dir := f.FilePathProxy.Dir(filePath)
	fileName := f.FilePathProxy.Base(filePath)
	// create hidden file path
	hiddenFilePath := f.FilePathProxy.Join(dir, "."+fileName)
	// rename file to hidden file
	if err := f.OsProxy.Rename(filePath, hiddenFilePath); err != nil {
		return -1, err
	}
	// append hidden file to hidden files
	f.HiddenFiles = append(f.HiddenFiles, hiddenFilePath)
	// return hidden file's index
	return len(f.HiddenFiles) - 1, nil
}

// RestoreFile restores hidden file to file.
func (f *FileHider) RestoreFile(index int) error {
	// get hidden file path
	hiddenFilePath := f.HiddenFiles[index]
	// check if file exists
	if _, err := f.OsProxy.Stat(hiddenFilePath); f.OsProxy.IsNotExist(err) {
		return err
	}
	// remove the leading dot from the file name
	restoredFilePath := f.FilePathProxy.Join(
		f.FilePathProxy.Dir(hiddenFilePath),
		f.StringsProxy.TrimPrefix(f.FilePathProxy.Base(hiddenFilePath), "."),
	)
	// rename hidden file to file
	if err := f.OsProxy.Rename(hiddenFilePath, restoredFilePath); err != nil {
		return err
	}
	// remove hidden file from hidden files
	f.HiddenFiles = append(f.HiddenFiles[:index], f.HiddenFiles[index+1:]...)
	return nil
}
