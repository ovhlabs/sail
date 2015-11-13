package operation

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(cmdOperationList)
	Cmd.AddCommand(cmdOperationAttach)
}

// Cmd service
var Cmd = &cobra.Command{
	Use:     "operation",
	Short:   "Operation commands: sail operation --help",
	Long:    `Operation commands: sail operation <command>`,
	Aliases: []string{"operations"},
}
