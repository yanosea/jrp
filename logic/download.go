package logic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yanosea/jrp/constant"
)

func Download(e Env, u User, fs FileSystem, hc HttpClient, io IO, gz Gzip) error {
	// get db file directory path
	dbFileDirPath, err := GetDBFileDirPath(e, u)
	if err != nil {
		return err
	}

	// if db file directory does not exist, create it
	if _, err := os.Stat(dbFileDirPath); os.IsNotExist(err) {
		os.MkdirAll(dbFileDirPath, 0755)
	}

	// if db file does not exist, download it
	dbFilePath := filepath.Join(dbFileDirPath, constant.WNJPN_DB_FILE_NAME)
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		// download db archive file
		resp, err := hc.Get(constant.WNJPN_DB_ARCHIVE_FILE_URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// save db archive file to temporary file
		tempFilePath := filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
		out, err := fs.Create(tempFilePath)
		if err != nil {
			return err
		}
		defer out.Close()
		if _, err := io.Copy(out, resp.Body); err != nil {
			return err
		}
		out.Seek(0, 0)

		// decompress db archive file to db file
		gz, err := gz.NewReader(out)
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

		// remove temporary file
		if err := fs.Remove(tempFilePath); err != nil {
			return err
		}
		// if db file is downloaded, print message
		fmt.Println(constant.DOWNLOAD_MESSAGE_SUCCEEDED)
	} else {
		// if db file already exists, print message
		fmt.Println(constant.DOWNLOAD_MESSAGE_ALREADY_DOWNLOADED)
	}

	return nil
}
