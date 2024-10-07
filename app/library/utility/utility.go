package utility

import (
	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/io"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/strconv"
)

// UtilityInterface is an interface for Utility.
type UtilityInterface interface {
	PrintlnWithWriter(writer ioproxy.WriterInstanceInterface, a ...any)
	GetMaxConvertibleString(args []string, def string) string
	GetLargerNumber(num int, argNum string) int
	CreateDirIfNotExist(dirPath string) error
}

// Utility is a struct that implements UtilityInterface.
type Utility struct {
	FmtProxy     fmtproxy.Fmt
	OsProxy      osproxy.Os
	StrconvProxy strconvproxy.Strconv
}

// New is a constructor for Utility.
func New(
	fmtProxy fmtproxy.Fmt,
	osProxy osproxy.Os,
	strconvProxy strconvproxy.Strconv,
) *Utility {
	return &Utility{
		FmtProxy:     fmtProxy,
		OsProxy:      osProxy,
		StrconvProxy: strconvProxy,
	}
}

// PrintlnWithWriter prints any with a writer.
func (u *Utility) PrintlnWithWriter(writer ioproxy.WriterInstanceInterface, a ...any) {
	u.FmtProxy.Fprintf(writer, u.FmtProxy.Sprintf("%s", a[0])+"\n")
}

// GetMaxConvertibleString gets the maximum number from args and converts it to a string.
func (u *Utility) GetMaxConvertibleString(args []string, def string) string {
	var maxArg string
	var maxValue int
	initialized := false

	for _, arg := range args {
		if convertedArg, err := u.StrconvProxy.Atoi(arg); err == nil {
			if !initialized || convertedArg > maxValue {
				// if the value is the first one or the value is larger than the max value
				maxValue = convertedArg
				maxArg = arg
				initialized = true
			}
		}
	}

	if initialized {
		// if there is less than 1 convertible arg, return the max arg
		return maxArg
	}

	// if there is no convertible args, return default value
	return def
}

// GetLargerNumber gets the larger number between num and argNum.
func (u *Utility) GetLargerNumber(num int, argNum string) int {
	if num <= 0 {
		// if num is less than 1, set num to 1
		num = 1
	}

	convertedArgNum, err := u.StrconvProxy.Atoi(argNum)
	if err != nil {
		// if argNum is not convertible, set argNum to 1
		convertedArgNum = 1
	}
	if convertedArgNum <= 0 {
		// if argNum is less than 1, set argNum to 1
		convertedArgNum = 1
	}

	// return the larger number
	if convertedArgNum > num {
		return convertedArgNum
	} else {
		return num
	}
}

// CreateDirIfNotExist creates a directory if it does not exist.
func (u *Utility) CreateDirIfNotExist(dirPath string) error {
	if _, err := u.OsProxy.Stat(dirPath); u.OsProxy.IsNotExist(err) {
		// if not exist, create dir
		return u.OsProxy.MkdirAll(dirPath, u.OsProxy.FileMode(0755))
	}
	return nil
}
