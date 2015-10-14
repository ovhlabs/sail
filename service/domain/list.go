package domain

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

var usageList = "Invalid usage. sail service domain list [[<application-name>/]<service-name>]. Please see sail domain list --help"
var domainHeadersDone = false

var cmdDomainList = &cobra.Command{
	Use:     "list",
	Short:   "List domains on the HTTP load balancer: sail service domain list [<application-name>[/<service-id>]]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Fprintln(os.Stderr, usageList)
			os.Exit(1)
		}

		namespace := ""
		service := ""

		// Parse namespace / service
		if len(args) >= 1 {
			t := strings.Split(args[0], "/")
			namespace = t[0]
			if len(t) >= 2 {
				service = t[1]
			} else if len(t) > 2 {
				fmt.Fprintln(os.Stderr, usageList)
				return
			}
		}

		domainList(namespace, service)
	},
}

func domainList(namespace, service string) {
	var apps []string

	// TODO: rewrite whithout the m(n+1)+1... (needs API)
	if len(namespace) > 0 {
		apps = append(apps, namespace)
	} else {
		apps = internal.GetListApplications(nil)
	}

	for _, namespace := range apps {
		domainListNamespace(namespace, service)
	}
}

func domainListNamespace(namespace, service string) {
	var services []string

	// TODO: rewrite whithout the m(n+1)+1... (needs API)
	if len(service) > 0 {
		services = append(services, service)
	} else {
		b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/services", namespace), nil)
		internal.Check(json.Unmarshal(b, &services))
	}

	for _, service := range services {
		domainListService(namespace, service)
	}
}

func domainListService(namespace, service string) {
	b := internal.ReqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/services/%s/attached-routes", namespace, service), nil)
	internal.FormatOutput(b, domainListFormatter)
}

func domainListFormatter(data []byte) {
	var routes []map[string]interface{}
	internal.Check(json.Unmarshal(data, &routes))

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)

	// below this: horrible hack. Do I feel ashamed: Yes.
	if !domainHeadersDone {
		titles := []string{"APP", "SERVICE", "DOMAIN", "METHOD", "PATTERN"}
		fmt.Fprintln(w, strings.Join(titles, "\t"))
		domainHeadersDone = true
	}

	for _, route := range routes {
		app := route["namespace"]
		service := route["service"]

		if app == nil {
			app = "-"
		}
		if service == nil {
			service = "-"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", app, service, route["domain"], route["method"], route["pattern"])
		w.Flush()
	}
}
