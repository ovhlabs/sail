package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdServiceList = &cobra.Command{
	Use:     "list",
	Short:   "List the docker services: sail service list [applicationName]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		serviceList(internal.GetListApplications(args))
	},
}

func serviceList(apps []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	titles := []string{"NAME", "REPOSITORY", "IMAGE ID", "STATE", "CONTAINERS", "CREATED", "NETWORK"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	services := []string{}
	var service map[string]interface{}
	for _, app := range apps {
		// Sanity checks
		err := internal.CheckName(app)
		internal.Check(err)

		b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/services", app), nil)
		internal.Check(json.Unmarshal(b, &services))
		sort.StringSlice(services).Sort()
		for _, serviceID := range services {
			b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/services/%s", app, serviceID), nil)
			internal.Check(json.Unmarshal(b, &service))

			ips := []string{}
			for _, container := range service["containers"].(map[string]interface{}) {
				for name, network := range container.(map[string]interface{})["network"].(map[string]interface{}) {
					ips = append(ips, fmt.Sprintf("%s:%s", name, network.(map[string]interface{})["ip"]))
				}
			}

			fmt.Fprintf(w, "%s/%s\t%s@%s\t%s\t%s\t%d\t%s\t%s\n",
				app, service["name"],
				service["repository"],
				service["repository_tag"],
				service["image"].(string)[:12],
				strings.ToUpper(service["state"].(string)),
				int(service["container_number"].(float64)),
				service["creation_date"].(string)[:19],
				strings.Join(ips, ","))
		}
	}
	w.Flush()
}
