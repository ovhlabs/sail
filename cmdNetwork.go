package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"
)

func init() {
	cmdNetwork.AddCommand(cmdNetworkList)
	cmdNetwork.AddCommand(cmdNetworkInspect)

	// TODO
	// sail networks add        Add a new private network
	// sail networks range-add  Add an allocation range to a private network
	// sail networks rm         Delete a private network

}

var cmdNetwork = &cobra.Command{
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
		networkList(getListApplications(args))
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

		b := reqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks/%s", t[0], t[1]), nil)
		check(json.Unmarshal(b, &network))

		brange := reqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks/%s/ranges", t[0], t[1]), nil)
		check(json.Unmarshal(brange, &ranges))

		network["range"] = ranges
		n, err := json.Marshal(network)
		check(err)
		fmt.Println(getJSON(n))
	}
}

func networkList(apps []string) {
	w := tabwriter.NewWriter(os.Stdout, 30, 1, 3, ' ', 0)
	titles := []string{"NAME", "SUBNET"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	networks := []string{}
	var network map[string]interface{}
	for _, app := range apps {
		b := reqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks", app), nil)
		check(json.Unmarshal(b, &networks))
		for _, networkID := range networks {
			b := reqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/networks/%s", app, networkID), nil)
			check(json.Unmarshal(b, &network))

			subnet := network["subnet"]
			if network["subnet"] == nil || network["subnet"] == "" {
				subnet = "-"
			}

			fmt.Fprintf(w, "%s\t%s\n", network["name"], subnet)
			w.Flush()
		}
	}
}
