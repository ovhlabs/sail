package container

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/runabove/sail/internal"
)

var cmdContainerAttach = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a container console: sail container attach <applicationName>/<containerId>",
	Long: `Attach to a container console: sail container attach <applicationName>/<containerId>
	"example: sail container attach my-app/myContainerId"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail container attach <applicationName>/<containerId>. Please see sail container attach --help")
		} else {
			containerAttach(args[0])
		}
	},
}

func containerAttach(containerID string) {
	t := strings.Split(containerID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid usage. sail container attach <applicationName>/<containerId>. Please see sail container attach --help")
	} else {
		internal.StreamPrint("GET", fmt.Sprintf("/applications/%s/containers/%s/attach", t[0], t[1]), nil)
		internal.ExitAfterCtrlC()
	}
}
