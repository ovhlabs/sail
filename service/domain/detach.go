package domain

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

var usageDomainDetach = "Invalid usage. sail service domain detach [<applicationName>/]<serviceId> <domain> <pattern> <method>. Please see sail service domain detach --help"
var cmdDomainDetach = &cobra.Command{
	Use:     "detach",
	Aliases: []string{"delete", "del", "rm", "remove"},
	Short:   "Detach a domain on the HTTP load balancer: sail service domain detach [<applicationName>/]<serviceId> <domain> <pattern> <method>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Fprintln(os.Stderr, usageDomainDetach)
		} else {
			serviceDomainDetach(args[0], args[1], domainStruct{Pattern: args[2], Method: args[3]})
		}
	},
}

func serviceDomainDetach(serviceID, domain string, args domainStruct) {
	// Split namespace and service
	host, app, service, tag, err := internal.ParseResourceName(serviceID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	} else if len(tag) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid service name. Please see sail service domain detach --help\n")
		os.Exit(1)
	}

	body, err := json.Marshal(args)
	internal.Check(err)

	// Sanity checks
	err = internal.CheckName(domain)
	internal.Check(err)

	path := fmt.Sprintf("/applications/%s/services/%s/attached-routes/%s", app, service, domain)
	data := internal.DeleteBodyWantJSON(path, body)

	internal.FormatOutput(data, func(data []byte) {
		fmt.Fprintf(os.Stderr, "Detached route %s %s%s from service %s/%s\n", args.Method, domain, args.Pattern, app, service)
	})
}
