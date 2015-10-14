package network

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

const cmdNetRangeAddUsage = "Add an allocation range to a private network: sail network range-add [<applicationName>/]<networkId> <ipFrom> <ipTo>"

var cmdNetworkRangeAdd = &cobra.Command{
	Use:     "rangeAdd",
	Short:   cmdNetRangeAddUsage,
	Long:    cmdNetRangeAddUsage,
	Aliases: []string{"range-add"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Fprintln(os.Stderr, cmdNetRangeAddUsage)
		} else {
			networkRangeAdd(args[0], args[1], args[2])
		}
	},
}

func networkRangeAdd(networkID, ipFrom, ipTo string) {
	// Split namespace and repository
	host, app, net, tag, err := internal.ParseResourceName(networkID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	} else if len(tag) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid network name. Please see sail network show --help\n")
		os.Exit(1)
	}

	path := fmt.Sprintf("/applications/%s/networks/%s/ranges/%s-%s", app, net, ipFrom, ipTo)
	internal.FormatOutputDef(internal.PostWantJSON(path))

}
