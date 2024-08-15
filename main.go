package main

import (
	"os"

	"github.com/yanosea/jrp/cmd"
)

var osExit = os.Exit

func main() {
	g := cmd.NewGlobalOption(os.Stdout, os.Stderr, os.Args[1:])
	osExit(g.Execute())
}
