package util

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

func FormatIndent(m string) string {
	return "  " + m
}

func PrintlnWithWriter(w io.Writer, a ...any) {
	fmt.Fprintln(w, fmt.Sprintf("%s", a[0]))
}

func PrintWithWriterWithBlankLineBelow(w io.Writer, a ...any) {
	fmt.Fprintln(w, fmt.Sprintf("%s\n", a[0]))
}

func PrintWithWriterWithBlankLineAbove(w io.Writer, a ...any) {
	fmt.Fprintln(w, fmt.Sprintf("\n%s", a[0]))
}

func PrintWithWriterBetweenBlankLine(w io.Writer, a ...any) {
	fmt.Fprintln(w, fmt.Sprintf("\n%s\n", a[0]))
}

func GetDBFileDirPath() (string, error) {
	// check if JRP_ENV is set
	var dbFileDir = os.Getenv(JRP_ENV)
	if dbFileDir == "" {
		// get current user
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		// default path ($XDG_DATA_HOME/jrp)
		dbFileDir = filepath.Join(user.HomeDir, ".local", "share", "jrp")
	}

	return filepath.Join(dbFileDir, WNJPN_DB_FILE_NAME), nil
}
