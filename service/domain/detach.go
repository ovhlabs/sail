package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

var usageDomainDetach = "Invalid usage. sail service detach <applicationName>/<serviceId> <domain> <pattern> <method>. Please see sail service detach --help"
var cmdDomainDetach = &cobra.Command{
	Use:     "detach",
	Aliases: []string{"delete", "del", "rm", "remove"},
	Short:   "Detach a domain on the HTTP load balancer: sail service domain detach <applicationName>/<serviceId> <domain> <pattern> <method>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Fprintln(os.Stderr, usageDomainDetach)
		} else {
			serviceDomainDetach(args[0], args[1], domainStruct{Pattern: args[2], Method: args[3]})
		}
	},
}

func serviceDomainDetach(serviceID, domain string, args domainStruct) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, usageDomainDetach)
		return
	}

	body, err := json.Marshal(args)
	internal.Check(err)

	path := fmt.Sprintf("/applications/%s/services/%s/attached-routes/%s", t[0], t[1], domain)
	internal.FormatOutputDef(internal.DeleteBodyWantJSON(path, body))

}
