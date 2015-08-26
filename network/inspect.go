package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sail/internal"
)

var cmdNetworkInspect = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect the docker private networks : sail network inspect <applicationName>/<networkId>",
	Long: `Inspect the docker private networks : sail network inspect <applicationName>/<networkId>
	\"example : sail network inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid usage. sail network inspect <applicationName>/<networkId>. Please see sail network inspect --help")
		} else {
			networkInspect(args[0])
		}
	},
}

func networkInspect(networkID string) {
	t := strings.Split(networkID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sail network inspect <applicationName>/<networkId>. Please see sail network inspect --help")
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
		fmt.Println(internal.GetJSON(n))
	}
}
