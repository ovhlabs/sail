package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

const cmdNetShowUsage = "Show the docker private networks: sail network show <applicationName>/<networkId>"

var cmdNetworkShow = &cobra.Command{
	Use:     "show",
	Aliases: []string{"inspect"},
	Short:   cmdNetShowUsage,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail network show <applicationName>/<networkId>. Please see sail network show --help")
		} else {
			networkShow(args[0])
		}
	},
}

func networkShow(networkID string) {
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

	var network map[string]interface{}
	var ranges []string

	b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks/%s", app, net), nil)
	internal.Check(json.Unmarshal(b, &network))

	brange := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks/%s/ranges", app, net), nil)
	internal.Check(json.Unmarshal(brange, &ranges))

	network["range"] = ranges
	n, err := json.Marshal(network)
	internal.Check(err)

	internal.FormatOutputDef(n)
}
