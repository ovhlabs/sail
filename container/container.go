package container

import (
	"fmt"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sail/internal"
)

func init() {
	Cmd.AddCommand(cmdContainerList)
	Cmd.AddCommand(cmdContainerInspect)
	Cmd.AddCommand(cmdContainerAttach)
	Cmd.AddCommand(cmdContainerLogs())
}

// Cmd container
var Cmd = &cobra.Command{
	Use:     "container",
	Short:   "Container commands : sail container --help",
	Long:    `Container commands : sail container <command>`,
	Aliases: []string{"c", "containers"},
}

var cmdContainerInspect = &cobra.Command{
	Use:     "inspect",
	Aliases: []string{"show"},
	Short:   "Inspect a docker container : sail container inspect <applicationName> <containerId>",
	Long: `Inspect a docker container : sail container inspect <applicationName> <containerId>
	\"example : sail container inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Invalid usage. sail container inspect <applicationName> <containerId>. Please see sail container inspect --help")
		} else {
			internal.FormatOutputDef(internal.GetWantJSON(fmt.Sprintf("/applications/%s/containers/%s", args[0], args[1])))
		}
	},
}
