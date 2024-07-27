package main

import (
	"os"

	"github.com/yanosea/jrp/cmd"
)

var exitFunc = os.Exit

func main() {
	exitFunc(cmd.Execute())
}
