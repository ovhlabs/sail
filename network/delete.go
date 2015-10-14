package network

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

const cmdNetDelUsage = "Remove a private network: sail network delete [<applicationName>/]<networkId>"

var cmdNetworkDelete = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm", "remove"},
	Short:   cmdNetDelUsage,
	Long:    cmdNetDelUsage,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail network delete <applicationName>/<networkId>. Please see sail network delete --help")
		} else {
			networkRemove(args[0])
		}
	},
}

func networkRemove(networkID string) {
	// Split namespace and repository
	host, app, net, tag, err := internal.ParseResourceName(networkID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	} else if len(tag) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid network name. Please see sail network delete --help\n")
		os.Exit(1)
	}

	path := fmt.Sprintf("/applications/%s/networks/%s", app, net)
	internal.FormatOutputDef(internal.DeleteWantJSON(path))
}
