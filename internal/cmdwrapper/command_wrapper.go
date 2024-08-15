package cmdwrapper

import (
	"github.com/spf13/cobra"
)

type ICommand interface {
	Execute() error
}

type CommandWrapper struct {
	Cmd *cobra.Command
}

func (cw *CommandWrapper) Execute() error {
	return cw.Cmd.Execute()
}

func NewCommandWrapper(cmd *cobra.Command) ICommand {
	return &CommandWrapper{Cmd: cmd}
}
