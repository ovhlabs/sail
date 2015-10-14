package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var domainHeadersDone = false

var cmdApplicationDomain = &cobra.Command{
	Use:     "domain",
	Short:   "Application Domain commands: sail application domain --help",
	Long:    `Application Domain commands: sail application domain <command>`,
	Aliases: []string{"domains"},
}

var cmdApplicationDomainList = &cobra.Command{
	Use:     "list",
	Short:   "List domains and routes on the HTTP load balancer: sail application domain list [applicationName]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		app := ""

		if len(args) == 1 && args[0] != "" {
			app = args[0]
		} else if len(args) > 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application domain list --help")
			return
		}

		domainList(app)
	},
}

var cmdApplicationDomainDetach = &cobra.Command{
	Use:     "detach",
	Short:   "Detach a domain from the HTTP load balancer: sail application domain detach <applicationName> <domainName>",
	Aliases: []string{"add"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application domain attach --help")
		} else {
			path := fmt.Sprintf("/applications/%s/attached-domains/%s", args[0], args[1])
			data := internal.DeleteWantJSON(path)

			internal.FormatOutput(data, func(data []byte) {
				fmt.Fprintf(os.Stderr, "Detached domain %s from application %s\n", args[1], args[0])
			})
		}
	},
}

func domainList(app string) {
	var apps []string

	// TODO: rewrite whithout the n+1... (needs API)
	if len(app) > 0 {
		apps = append(apps, app)
	} else {
		apps = internal.GetListApplications(nil)
	}

	for _, app := range apps {
		domainListApplication(app)
	}
}

func domainListApplication(app string) {
	b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/attached-domains", app), nil)
	internal.FormatOutput(b, domainListFormatter)
}

func domainListFormatter(data []byte) {
	var domains map[string][]map[string]interface{}
	internal.Check(json.Unmarshal(data, &domains))

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)

	// below this: horrible hack. Do I feel ashamed: Yes.
	if !domainHeadersDone {
		titles := []string{"APP", "SERVICE", "DOMAIN", "METHOD", "PATTERN"}
		fmt.Fprintln(w, strings.Join(titles, "\t"))
		domainHeadersDone = true
	}

	for domain, routes := range domains {
		for _, route := range routes {
			service := route["service"]
			app := route["namespace"]

			if app == nil {
				app = "-"
			}
			if service == nil {
				service = "-"
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", app, service, domain, route["method"], route["pattern"])
			w.Flush()
		}
	}
}
