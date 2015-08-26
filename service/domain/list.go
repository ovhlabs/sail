package domain

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

var usageList = "Invalid usage. sail service domain list <applicationName>/<networkId>. Please see sail domain list --help"

var cmdDomainList = &cobra.Command{
	Use:     "list",
	Short:   "List domains on the HTTP load balancer : sail service domain list <applicationName>/<serviceId>",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(usageList)
		} else {
			domainList(args[0])
		}
	},
}

func domainList(serviceID string) {

	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Println(usageList)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 30, 1, 3, ' ', 0)
	titles := []string{"DOMAIN", "METHOD", "PATTERN"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	var routes []map[string]interface{}

	b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/services/%s/attached-routes", t[0], t[1]), nil)
	internal.Check(json.Unmarshal(b, &routes))
	for _, route := range routes {
		fmt.Fprintf(w, "%s\t%s\t%s\n", route["domain"], route["method"], route["pattern"])
		w.Flush()
	}
}
