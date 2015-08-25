package container

import (
	"fmt"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sailgo/internal"
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
	Short:   "Container commands : sailgo container --help",
	Long:    `Container commands : sailgo container <command>`,
	Aliases: []string{"c", "containers"},
}

var cmdContainerInspect = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a docker container : sailgo container inspect <applicationName> <containerId>",
	Long: `Inspect a docker container : sailgo container inspect <applicationName> <containerId>
	\"example : sailgo container inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Invalid usage. sailgo container inspect <applicationName> <containerId>. Please see sailgo container inspect --help")
		} else {
			fmt.Println(internal.GetWantJSON(fmt.Sprintf("/applications/%s/containers/%s", args[0], args[1])))
		}
	},
}
