package container

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sailgo/internal"
)

var cmdContainerAttach = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a container console : sailgo container attach <applicationName>/<containerId>",
	Long: `Attach to a container console : sailgo container attach <applicationName>/<containerId>
	\"example : sailgo container attach myApp myContainerId"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid usage. sailgo container attach <applicationName>/<containerId>. Please see sailgo container attach --help")
		} else {
			containerAttach(args[0])
		}
	},
}

func containerAttach(containerID string) {
	t := strings.Split(containerID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sailgo service inspect <applicationName>/<serviceId>. Please see sailgo service inspect --help")
	} else {
		internal.StreamWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/containers/%s/attach", t[0], t[1]), nil)
	}
}
