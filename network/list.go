package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sail/internal"
)

var cmdNetworkList = &cobra.Command{
	Use:     "list",
	Short:   "List the docker private networks : sail network list [applicationName]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		networkList(internal.GetListApplications(args))
	},
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
