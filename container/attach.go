package container

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdContainerAttach = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a container console: sail container attach [<applicationName>/]<containerId>",
	Long: `Attach to a container console: sail container attach [<applicationName>/]<containerId>
	"example: sail container attach my-app/myContainerId"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail container attach [<applicationName>/]<containerId>. Please see sail container attach --help")
		} else {
			containerAttach(args[0])
		}
	},
}

func containerAttach(containerID string) {
	// Split namespace and container
	host, app, container, tag, err := internal.ParseResourceName(containerID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	} else if len(tag) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid container name. Please see sail container attach --help\n")
		os.Exit(1)
	}

	internal.StreamPrint("GET", fmt.Sprintf("/applications/%s/containers/%s/attach", app, container), nil)
	internal.ExitAfterCtrlC()
}
