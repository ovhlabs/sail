package network

import (
	"fmt"
	"os"
	"strings"

	"stash.ovh.net/sailabove/sail/internal"

	"github.com/spf13/cobra"
)

var cmdNetworkRm = &cobra.Command{
	Use:     "rm",
	Short:   "Remove a private network: sail network rm <applicationName>/<networkId>",
	Long:    `Remove a private network: sail network rm <applicationName>/<networkId>`,
	Aliases: []string{"delete", "remove"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail network rm <applicationName>/<networkId>. Please see sail network rm --help")
		} else {
			networkRemove(args[0])
		}
	},
}

func networkRemove(networkID string) {
	t := strings.Split(networkID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid usage. sail network rm <applicationName>/<networkId>. Please see sail network rm --help")
		return
	}

	path := fmt.Sprintf("/applications/%s/networks/%s", t[0], t[1])
	internal.FormatOutputDef(internal.DeleteWantJSON(path))
}
