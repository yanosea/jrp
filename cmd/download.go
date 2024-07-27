package cmd

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/constant"
	"github.com/yanosea/jrp/util"
)

type downloadOption struct {
	Out    io.Writer
	ErrOut io.Writer
}

func newDownloadCommand(globalOption *GlobalOption) *cobra.Command {
	o := &downloadOption{}
	cmd := &cobra.Command{
		Use:   constant.DOWNLOAD_USE,
		Short: constant.DOWNLOAD_SHORT,
		Long:  constant.DOWNLOAD_LONG,
		RunE: func(cmd *cobra.Command, args []string) error {

			o.Out = globalOption.Out
			o.ErrOut = globalOption.ErrOut

			return o.download()
		},
	}

	o.Out = globalOption.Out
	o.ErrOut = globalOption.ErrOut
	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)

	cmd.SetHelpTemplate(constant.DOWNLOAD_HELP_TEMPLATE)

	return cmd
}

func (o *downloadOption) download() error {
	// get the directory of wnjpn.db from environment
	dbFileDirPath, err := util.GetDBFileDirPath()
	if err != nil {
		return err
	}

	// create the directory if it doesn't exist
	if _, err := os.Stat(dbFileDirPath); os.IsNotExist(err) {
		os.MkdirAll(dbFileDirPath, 0755)
	}

	// download the database file if it doesn't exist
	var dbFilePath = filepath.Join(dbFileDirPath, constant.WNJPN_DB_FILE_NAME)
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		resp, err := http.Get(constant.WNJPN_DB_ARCHIVE_FILE_URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// save the downloaded file to a temporary file
		var tempFilePath = filepath.Join(os.TempDir(), constant.WNJPN_DB_ARCHIVE_FILE_NAME)
		out, err := os.Create(tempFilePath)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		_, err = out.Seek(0, 0)
		if err != nil {
			return err
		}

		// decompress the downloaded file
		r, err := gzip.NewReader(out)
		if err != nil {
			return err
		}
		defer r.Close()

		// save the decompressed file to the database file
		f, err := os.Create(dbFilePath)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, r)
		if err != nil {
			return err
		}

		// remove the temporary file
		err = os.Remove(tempFilePath)
		if err != nil {
			return err
		}
	}

	// already downloaded
	fmt.Println(constant.DOWNLOAD_MESSAGE_ALREADY_DOWNLOADED)

	return nil
}
