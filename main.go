package main

import (
	"os"

	"github.com/yanosea/jrp/cmd"
)

var osExit = os.Exit

func main() {
	osExit(cmd.Execute())
}
