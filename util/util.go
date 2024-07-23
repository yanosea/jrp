package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
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

func GetDBFileDirPath() string {
	var dbFileDir = os.Getenv(JRP_ENV)
	if dbFileDir == "" {
		// get home directory
		var homeDir string
		if runtime.GOOS == "windows" {
			homeDir = os.Getenv("USERPROFILE")
		} else {
			homeDir = os.Getenv("HOME")
		}
		// default path ($XDG_DATA_HOME/jrp)
		dbFileDir = filepath.Join(homeDir, ".local", "share", "jrp")
	}

	return dbFileDir
}
