package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sail/internal"
)

var cmdNetworkShow = &cobra.Command{
	Use:     "show",
	Aliases: []string{"inspect"},
	Short:   "Show the docker private networks: sail network show <applicationName>/<networkId>",
	Long: `Show the docker private networks: sail network show <applicationName>/<networkId>
	\"example: sail network show my-app"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail network show <applicationName>/<networkId>. Please see sail network show --help")
		} else {
			networkShow(args[0])
		}
	},
}

func networkShow(networkID string) {
	t := strings.Split(networkID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid usage. sail network show <applicationName>/<networkId>. Please see sail network show --help")
	} else {

		var network map[string]interface{}
		var ranges []string

		b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks/%s", t[0], t[1]), nil)
		internal.Check(json.Unmarshal(b, &network))

		brange := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks/%s/ranges", t[0], t[1]), nil)
		internal.Check(json.Unmarshal(brange, &ranges))

		network["range"] = ranges
		n, err := json.Marshal(network)
		internal.Check(err)
		internal.FormatOutputDef(n)
	}
}
