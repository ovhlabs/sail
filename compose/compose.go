package compose

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(cmdComposeUp)
	Cmd.AddCommand(cmdComposeGet)
}

// Cmd compose
var Cmd = &cobra.Command{
	Use:     "compose",
	Short:   "Compose commands: sail compose --help",
	Long:    `Compose commands: sail compose <command>`,
	Aliases: []string{"comp"},
}

var cmdComposeUp = &cobra.Command{
	Use:   "up",
	Short: "sail compose up <namespace>",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO sail compose up
		fmt.Fprintln(os.Stderr, "sail compose up TO BE IMPLEMENTED")
	},
}

var cmdComposeGet = &cobra.Command{
	Use:   "get",
	Short: "sail compose get <namespace>",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO  sail compose get
		fmt.Fprintln(os.Stderr, "sail compose get TO BE IMPLEMENTED")
	},
}
