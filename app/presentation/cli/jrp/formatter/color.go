package formatter

import (
	"github.com/fatih/color"
)

var (
	// Blue is a proxy of fatih/color.Blue.
	Blue = color.New(color.FgBlue).SprintFunc()
	// Green is a proxy of fatih/color.Green.
	Green = color.New(color.FgGreen).SprintFunc()
	// Red is a proxy of fatih/color.Red.
	Red = color.New(color.FgRed).SprintFunc()
	// Yellow is a proxy of fatih/color.Yellow.
	Yellow = color.New(color.FgYellow).SprintFunc()
)
