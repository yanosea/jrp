package util

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/yanosea/jrp/constant"
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

func GetDBFileDirPath(provider UserProvider) (string, error) {
	// check if JRP_ENV is set
	var dbFileDirPath = os.Getenv(constant.JRP_ENV_WORDNETJP_DIR)
	if dbFileDirPath == "" {
		// get current user
		user, err := provider.Current()
		if err != nil {
			return "", err
		}
		// default path ($XDG_DATA_HOME/jrp)
		dbFileDirPath = filepath.Join(user.HomeDir, ".local", "share", "jrp")
	}

	return dbFileDirPath, nil
}

// for testing
type UserProvider interface {
	Current() (*user.User, error)
}

type DefaultUserProvider struct{}

func (d DefaultUserProvider) Current() (*user.User, error) {
	return user.Current()
}
