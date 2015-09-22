package container

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/runabove/sail/internal"
)

var cmdContainerList = &cobra.Command{
	Use:     "list",
	Short:   "List docker containers: sail container list [applicationName]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		containerList(internal.GetListApplications(args))
	},
}

func containerList(apps []string) {
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	titles := []string{"APPLICATION", "SERVICE", "CONTAINER", "STATE", "DEPLOYED"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	containers := []string{}
	var container map[string]interface{}
	for _, app := range apps {
		b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/containers", app), nil)
		internal.Check(json.Unmarshal(b, &containers))
		for _, containerID := range containers {
			b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/containers/%s", app, containerID), nil)
			internal.Check(json.Unmarshal(b, &container))
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", app, container["service"], container["name"], strings.ToUpper(container["state"].(string)), container["deployment_date"])
			w.Flush()
		}
	}
}
