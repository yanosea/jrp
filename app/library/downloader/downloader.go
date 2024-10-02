package downloader

import (
	"github.com/yanosea/jrp/app/library/utility"
	"github.com/yanosea/jrp/app/proxy/filepath"
	"github.com/yanosea/jrp/app/proxy/gzip"
	"github.com/yanosea/jrp/app/proxy/http"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
)

// Downloadable is an interface that provides a function to download wnjpn db file.
type Downloadable interface {
	DownloadWNJpnDBFile(wnJpnDBFileDirPath string) (DownloadStatus, error)
}

// Downloader is a struct that implements Downloadable interface.
type Downloader struct {
	FilepathProxy filepathproxy.FilePath
	GzipProxy     gzipproxy.Gzip
	HttpProxy     httpproxy.Http
	IoProxy       ioproxy.Io
	OsProxy       osproxy.Os
	Utility       utility.UtilityInterface
}

// New is a constructor of Downloader.
func New(
	filepathProxy filepathproxy.FilePath,
	gzipProxy gzipproxy.Gzip,
	httpProxy httpproxy.Http,
	ioProxy ioproxy.Io,
	osProxy osproxy.Os,
	utility utility.UtilityInterface,
) *Downloader {
	return &Downloader{
		FilepathProxy: filepathProxy,
		GzipProxy:     gzipProxy,
		HttpProxy:     httpProxy,
		IoProxy:       ioProxy,
		OsProxy:       osProxy,
		Utility:       utility,
	}
}

// DownloadWNJpnDBFile downloads wnjpn db file.
func (d *Downloader) DownloadWNJpnDBFile(wnJpnDBFileDirPath string) (DownloadStatus, error) {
	// create dir if not exist
	if err := d.Utility.CreateDirIfNotExist(wnJpnDBFileDirPath); err != nil {
		// if failed to create dir, return failure
		return DownloadedFailed, err
	}

	// check if db file is already downloaded
	dbFilePath := d.FilepathProxy.Join(wnJpnDBFileDirPath, WNJPN_DB_FILE_NAME)
	if _, err := d.OsProxy.Stat(dbFilePath); d.OsProxy.IsNotExist(err) {
		// if not downloaded, download and extract db file
		return d.downloadAndExtractDBFile(dbFilePath)
	}

	// if already downloaded, return
	return DownloadedAlready, nil
}

// downloadAndExtractDBFile downloads and extracts wnjapn db file.
func (d *Downloader) downloadAndExtractDBFile(dbFilePath string) (DownloadStatus, error) {
	// download gzip file
	resp, err := d.downloadGzipFile()
	if err != nil {
		return DownloadedFailed, err
	}
	defer resp.FieldResponse.Body.Close()

	// save to temp file
	tempFilePath, err := d.saveToTempFile(resp.FieldResponse.Body)
	if err != nil {
		return DownloadedFailed, err
	}
	defer d.OsProxy.Remove(tempFilePath)

	// extract gzip file
	if err := d.extractGzipFile(tempFilePath, dbFilePath); err != nil {
		return DownloadedFailed, err
	}

	return DownloadedSuccessfully, nil
}

// downloadGzipFile downloads gzip file.
func (d *Downloader) downloadGzipFile() (*httpproxy.ResponseInstance, error) {
	// download gzip file
	return d.HttpProxy.Get(WNJPN_DB_ARCHIVE_FILE_URL)
}

// saveToTempFile saves body to temp file.
func (d *Downloader) saveToTempFile(body ioproxy.ReaderInstanceInterface) (string, error) {
	// create temp file
	tempFilePath := d.FilepathProxy.Join(d.OsProxy.TempDir(), WNJPN_DB_ARCHIVE_FILE_NAME)
	out, err := d.OsProxy.Create(tempFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// copy downloaded file to temp file
	if _, err := d.IoProxy.Copy(out, body); err != nil {
		return "", err
	}

	// seek to start
	if _, err := out.Seek(0, ioproxy.SeekStart); err != nil {
		return "", err
	}

	return tempFilePath, nil
}

// extractGzipFile extracts gzip file.
func (d *Downloader) extractGzipFile(srcPath, destPath string) error {
	// open gzip file
	file, err := d.OsProxy.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()
	gz, err := d.GzipProxy.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()

	// create file to save
	out, err := d.OsProxy.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// copy gzip file to dest file
	if _, err := d.IoProxy.Copy(out, gz); err != nil {
		return err
	}

	return nil
}
