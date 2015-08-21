package repository

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

var cmdRepositoryList = &cobra.Command{
	Use:     "list",
	Short:   "List the docker repository : sailgo repository list [applicationName]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		repositoryList(internal.GetListApplications(args))
	},
}

func repositoryList(apps []string) {
	w := tabwriter.NewWriter(os.Stdout, 30, 1, 3, ' ', 0)
	titles := []string{"NAME", "TAG", "TYPE", "PRIVACY", "SOURCE"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	repositories := []string{}
	var repository map[string]interface{}
	for _, app := range apps {
		b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/repositories/%s", app), nil)
		internal.Check(json.Unmarshal(b, &repositories))
		for _, repositoryID := range repositories {
			b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/repositories/%s/%s", app, repositoryID), nil)
			internal.Check(json.Unmarshal(b, &repository))

			tags := repository["tags"]
			if tags == "" {
				tags = "-"
			}
			source := repository["source"]
			if source == nil || source == "" {
				source = "-"
			}
			fmt.Fprintf(w, "%s/%s\t%s\t%s\t%s\t%s\n", app, repository["name"], tags, repository["type"], repository["privacy"], source)
			w.Flush()
		}
	}
}
