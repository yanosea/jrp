// Package main is the entry point of jrp.
package main

import (
	"os"

	"github.com/yanosea/jrp/app/proxy/fmt"
	"github.com/yanosea/jrp/app/proxy/os"
	"github.com/yanosea/jrp/app/proxy/strconv"
	"github.com/yanosea/jrp/cmd"
)

// fmtProxy is variable for fmt.Proxy.
var fmtProxy = fmtproxy.New()

// osProxy is variable for os.Proxy.
var osProxy = osproxy.New()

// strconvProxy is variable for strconv.Proxy.
var strconvProxy = strconvproxy.New()

// osExit is variable for os.Exit.
var osExit = os.Exit

// main is the entry point of jrp.
func main() {
	g := cmd.NewGlobalOption(fmtProxy, osProxy, strconvProxy)
	osExit(g.Execute())
}
