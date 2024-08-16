package logic

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/internal/fs"
	"github.com/yanosea/jrp/internal/gzip"
	"github.com/yanosea/jrp/internal/httpclient"
	"github.com/yanosea/jrp/internal/iomanager"
	"github.com/yanosea/jrp/internal/spinnerservice"
	"github.com/yanosea/jrp/internal/usermanager"
)

type Downloader interface {
	Download() error
}

type DBFileDownloader struct {
	User       usermanager.UserProvider
	FileSystem fs.FileManager
	HttpClient httpclient.HTTPClient
	IO         iomanager.IOHelper
	Gzip       gzip.GzipHandler
	Spinner    spinnerservice.SpinnerService
}

func NewDBFileDownloader(u usermanager.UserProvider, f fs.FileManager,
	h httpclient.HTTPClient, i iomanager.IOHelper, g gzip.GzipHandler, s spinnerservice.SpinnerService) *DBFileDownloader {
	return &DBFileDownloader{
		User:       u,
		FileSystem: f,
		HttpClient: h,
		IO:         i,
		Gzip:       g,
		Spinner:    s,
	}
}

func (d *DBFileDownloader) Download() error {
	// create DBFileDirPathGetter instance
	dbFileDirPathGetter := NewDBFileDirPathGetter(d.User)

	// get db file directory path
	dbFileDirPath, err := dbFileDirPathGetter.GetFileDirPath()
	if err != nil {
		return err
	}

	// if db file directory does not exist, create it
	if _, err := os.Stat(dbFileDirPath); os.IsNotExist(err) {
		if err := d.FileSystem.MkdirAll(dbFileDirPath, os.FileMode(0755)); err != nil {
			return err
		}
	}

	// if db file does not exist, download it
	dbFilePath := filepath.Join(dbFileDirPath, constant.WNJPN_DB_FILE_NAME)
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		// spinner settings
		if err := d.Spinner.SetColor("yellow"); err != nil {
			return err
		}
		d.Spinner.SetSuffix(color.YellowString(constant.DOWNLOAD_MESSAGE_DOWNLOADING))
		// start spinner
		d.Spinner.Start()

		// download db archive file
		resp, err := d.HttpClient.Get(constant.WNJPN_DB_ARCHIVE_FILE_URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// save db archive file to temporary file
		tempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
		out, err := d.FileSystem.Create(tempFilePath)
		if err != nil {
			return err
		}
		defer out.Close()
		if _, err := d.IO.Copy(out, resp.Body); err != nil {
			return err
		}
		if _, err := out.Seek(0, io.SeekStart); err != nil {
			return err
		}

		// decompress db archive file to db file
		gz, err := d.Gzip.NewReader(out)
		if err != nil {
			return err
		}
		defer gz.Close()
		f, err := d.FileSystem.Create(dbFilePath)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := d.IO.Copy(f, gz); err != nil {
			return err
		}

		// remove temporary file
		if err := d.FileSystem.RemoveAll(tempFilePath); err != nil {
			return err
		}

		// stop spinner
		d.Spinner.Stop()

		// if db file is downloaded successfully, print message
		fmt.Println(color.GreenString(constant.DOWNLOAD_MESSAGE_SUCCEEDED))
	} else {
		// if db file already exists, print message
		fmt.Println(color.GreenString(constant.DOWNLOAD_MESSAGE_ALREADY_DOWNLOADED))
	}

	return nil
}
