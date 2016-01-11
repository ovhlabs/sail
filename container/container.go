package container

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/runabove/sail/internal"
)

func init() {
	Cmd.AddCommand(cmdContainerList)
	Cmd.AddCommand(cmdContainerShow)
	Cmd.AddCommand(cmdContainerAttach)
	Cmd.AddCommand(cmdContainerLogs())
}

// Cmd container
var Cmd = &cobra.Command{
	Use:     "container",
	Short:   "Container commands: sail container --help",
	Long:    `Container commands: sail container <command>`,
	Aliases: []string{"c", "containers"},
}

var cmdContainerShow = &cobra.Command{
	Use:     "show",
	Aliases: []string{"inspect"},
	Short:   "Show a docker container: sail container show <containerId>",
	Long: `Show a docker container: sail container show <containerId>
	\"example: sail container show my-app my-container"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var container string

		if len(args) == 1 {
			container = args[0]
		} else if len(args) == 2 {
			container = args[1]
		} else {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail container show <applicationName> <containerId>. Please see sail container show --help")
			os.Exit(1)
		}
		internal.FormatOutputDef(internal.GetWantJSON(fmt.Sprintf("/containers/%s", container)))
	},
}
