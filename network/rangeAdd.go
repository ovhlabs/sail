package network

import (
	"fmt"
	"strings"

	"stash.ovh.net/sailabove/sail/internal"

	"github.com/spf13/cobra"
)

var cmdNetworkRangeAdd = &cobra.Command{
	Use:     "rangeAdd",
	Short:   "Add an allocation range to a private network : sail network range-add <applicationName>/<networkId> <ipFrom> <ipTo>",
	Long:    `Add an allocation range to a private network : sail network range-add <applicationName>/<networkId> <ipFrom> <ipTo>`,
	Aliases: []string{"range-add"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("Invalid usage. sail network range-add <applicationName>/<networkId> <ipFrom> <ipTo>. Please see sail network range-add --help")
		} else {
			networkRangeAdd(args[0], args[1], args[2])
		}
	},
}

func networkRangeAdd(networkID, ipFrom, ipTo string) {
	t := strings.Split(networkID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sail network range-add <applicationName>/<networkId> <ipFrom> <ipTo>. Please see sail network range-add --help")
		return
	}

	path := fmt.Sprintf("/applications/%s/networks/%s/ranges/%s-%s", t[0], t[1], ipFrom, ipTo)
	fmt.Println(internal.PostWantJSON(path))

}
