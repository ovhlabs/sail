package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sailgo/internal"
)

func init() {
	Cmd.AddCommand(cmdNetworkList)
	Cmd.AddCommand(cmdNetworkInspect)

	// TODO
	// sail networks add        Add a new private network
	// sail networks range-add  Add an allocation range to a private network
	// sail networks rm         Delete a private network

}

var Cmd = &cobra.Command{
	Use:     "network",
	Short:   "Network commands : sailgo network --help",
	Long:    `Network commands : sailgo network <command>`,
	Aliases: []string{"networks"},
}

var cmdNetworkList = &cobra.Command{
	Use:     "list",
	Short:   "List the docker private networks : sailgo network list [applicationName]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		networkList(internal.GetListApplications(args))
	},
}

var cmdNetworkInspect = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect the docker private networks : sailgo network inspect <applicationName>/<networkId>",
	Long: `Inspect the docker private networks : sailgo network inspect <applicationName>/<networkId>
	\"example : sailgo network inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid usage. sailgo network inspect <applicationName>/<networkId>. Please see sailgo network inspect --help")
		} else {
			networkInspect(args[0])
		}
	},
}

func networkInspect(networkID string) {
	t := strings.Split(networkID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sailgo network inspect <applicationName>/<networkId>. Please see sailgo network inspect --help")
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

func networkList(apps []string) {
	w := tabwriter.NewWriter(os.Stdout, 30, 1, 3, ' ', 0)
	titles := []string{"NAME", "SUBNET"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	networks := []string{}
	var network map[string]interface{}
	for _, app := range apps {
		b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks", app), nil)
		internal.Check(json.Unmarshal(b, &networks))
		for _, networkID := range networks {
			b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks/%s", app, networkID), nil)
			internal.Check(json.Unmarshal(b, &network))

			subnet := network["subnet"]
			if network["subnet"] == nil || network["subnet"] == "" {
				subnet = "-"
			}

			fmt.Fprintf(w, "%s\t%s\n", network["name"], subnet)
			w.Flush()
		}
	}
}
