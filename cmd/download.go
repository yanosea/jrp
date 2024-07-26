package cmd

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/yanosea/jrp/util"
)

const (
	download_help_template = `📥 Download Japanese Wordnet and English WordNet in an sqlite3 database from the official site.

You have to download Japanese Wordnet and English WordNet in an sqlite3 database to use jrp at first.
jrp will download archive file from the official site and decompress it to the database file.

You can set the directory of the database file to the environment variable "JRP_WORDNETJP_DIR".
The default directory is "~/.local/share/jrp" ("$XDG_DATA_HOME/jrp").

Usage:
  jrp download [flags]

Flags:
  -h, --help   🤝 help for download
`
	download_use   = "download"
	download_short = "📥 Download Japanese Wordnet and English WordNet in an sqlite3 database from the official site."
	download_long  = `📥 Download Japanese Wordnet and English WordNet in an sqlite3 database from the official site.

You have to download Japanese Wordnet and English WordNet in an sqlite3 database to use jrp at first.
jrp will download archive file from the official site and decompress it to the database file.

You can set the directory of the database file to the environment variable "JRP_WORDNETJP_DIR".
The default directory is "$XDG_DATA_HOME/jrp".
`
	download_message_already_downloaded = "✅ You are ready to use jrp!"
)

type downloadOption struct {
	Out    io.Writer
	ErrOut io.Writer
}

func newDownloadCommand(globalOption *GlobalOption) *cobra.Command {
	o := &downloadOption{}
	cmd := &cobra.Command{
		Use:   download_use,
		Short: download_short,
		Long:  download_long,
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

	cmd.SetHelpTemplate(download_help_template)

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
	var dbFilePath = filepath.Join(dbFileDirPath, util.WNJPN_DB_FILE_NAME)
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		resp, err := http.Get(util.WNJPN_DB_ARCHIVE_FILE_URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// save the downloaded file to a temporary file
		var tempFilePath = filepath.Join(os.TempDir(), util.WNJPN_DB_ARCHIVE_FILE_NAME)
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

	return nil
}
