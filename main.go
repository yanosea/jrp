package main

import (
	"os"

	"github.com/yanosea/jrp/cmd"
)

var osExit = os.Exit

func main() {
	g := &cmd.GlobalOption{Out: os.Stdout, ErrOut: os.Stderr}
	osExit(g.Execute())
}
