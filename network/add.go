package network

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

const cmdNetAddUsage = "Add a new private network: sail network add [<applicationName>/]<networkId> subnet"

var cmdNetworkAdd = &cobra.Command{
	Use:   "add",
	Short: cmdNetAddUsage,
	Long:  cmdNetAddUsage,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, cmdNetAddUsage)
		} else {
			n := networkAddStruct{Subnet: args[1]}
			networkAdd(args[0], n)
		}
	},
}

type networkAddStruct struct {
	Subnet string `json:"subnet"`
}

func networkAdd(networkID string, args networkAddStruct) {
	// Split namespace and repository
	host, app, net, tag, err := internal.ParseResourceName(networkID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	} else if len(tag) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid network name. Please see sail network add --help\n")
		os.Exit(1)
	}

	body, err := json.Marshal(args)
	internal.Check(err)

	path := fmt.Sprintf("/applications/%s/networks/%s", app, net)
	internal.FormatOutputDef(internal.PostBodyWantJSON(path, body))

}
