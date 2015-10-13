package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

var usageDomainAttach = "Invalid usage. sail service domain attach <applicationName>/<serviceId> <domain> [<pattern> [<method>]]. Please see sail service domain attach --help"
var cmdDomainAttach = &cobra.Command{
	Use:     "attach",
	Short:   "Attach a domain on the HTTP load balancer: sail service domain attach <applicationName>/<serviceId> <domain> [<pattern> [<method>]]",
	Aliases: []string{"add"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 || len(args) > 4 {
			fmt.Fprintln(os.Stderr, usageDomainAttach)
			os.Exit(1)
		}

		pattern := "/"
		method := "*"
		if len(args) >= 3 {
			pattern = args[2]
		}
		if len(args) >= 4 {
			method = args[3]
		}

		serviceDomainAttach(args[0], args[1], pattern, method)
	},
}

func serviceDomainAttach(serviceID, domain, pattern, method string) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, usageDomainAttach)
		return
	}

	args := domainStruct{Pattern: pattern, Method: method}
	body, err := json.Marshal(args)
	internal.Check(err)

	path := fmt.Sprintf("/applications/%s/services/%s/attached-routes/%s", t[0], t[1], domain)
	internal.FormatOutputDef(internal.PostBodyWantJSON(path, body))

}
