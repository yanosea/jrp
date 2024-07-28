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

func Download() error {
	dbFileDirPath, err := util.GetDBFileDirPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dbFileDirPath); os.IsNotExist(err) {
		os.MkdirAll(dbFileDirPath, 0755)
	}

	dbFilePath := filepath.Join(dbFileDirPath, constant.WNJPN_DB_FILE_NAME)
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		resp, err := http.Get(constant.WNJPN_DB_ARCHIVE_FILE_URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		tempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
		out, err := os.Create(tempFilePath)
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

		f, err := os.Create(dbFilePath)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, gz); err != nil {
			return err
		}

		if err := os.Remove(tempFilePath); err != nil {
			return err
		}
		fmt.Println(constant.DOWNLOAD_MESSAGE_SUCCEEDED)
	} else {
		fmt.Println(constant.DOWNLOAD_MESSAGE_ALREADY_DOWNLOADED)
	}

	return nil
}
