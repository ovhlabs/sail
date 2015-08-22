package network

import (
	"fmt"
	"strings"

	"stash.ovh.net/sailabove/sailgo/internal"

	"github.com/spf13/cobra"
)

var cmdNetworkRm = &cobra.Command{
	Use:     "rm",
	Short:   "Remove a private network : sailgo network rm <applicationName>/<networkId>",
	Long:    `Remove a private network : sailgo network rm <applicationName>/<networkId>`,
	Aliases: []string{"delete", "remove"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid usage. sailgo network rm <applicationName>/<networkId>. Please see sailgo network rm --help")
		} else {
			networkRemove(args[0])
		}
	},
}

func networkRemove(networkID string) {
	t := strings.Split(networkID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sailgo network rm <applicationName>/<networkId>. Please see sailgo network rm --help")
		return
	}

	path := fmt.Sprintf("/applications/%s/networks/%s", t[0], t[1])
	fmt.Println(internal.DeleteWantJSON(path))
}
