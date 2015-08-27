package compose

import (
	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(cmdComposeUp())
	Cmd.AddCommand(cmdComposeGet())
}

// Cmd compose
var Cmd = &cobra.Command{
	Use:     "compose",
	Short:   "Compose commands: sail compose --help",
	Long:    `Compose commands: sail compose <command>`,
	Aliases: []string{"comp"},
}
