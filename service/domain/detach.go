package domain

import (
	"encoding/json"
	"fmt"
	"strings"

	"stash.ovh.net/sailabove/sailgo/internal"

	"github.com/spf13/cobra"
)

var usageDomainDetach = "Invalid usage. sailgo service detach <applicationName>/<serviceId> <domain> <pattern> <method>. Please see sailgo service detach --help"
var cmdDomainDetach = &cobra.Command{
	Use:     "detach",
	Short:   "Detach a domain on the HTTP load balancer : sailgo service domain detach <applicationName>/<serviceId> <domain> <pattern> <method>",
	Aliases: []string{"rm", "remove", "delete"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println(usageDomainDetach)
		} else {
			serviceDomainDetach(args[0], args[1], domainStruct{Pattern: args[2], Method: args[3]})
		}
	},
}

func serviceDomainDetach(serviceID, domain string, args domainStruct) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Println(usageDomainDetach)
		return
	}

	body, err := json.Marshal(args)
	internal.Check(err)

	path := fmt.Sprintf("/applications/%s/services/%s/attached-routes/%s", t[0], t[1], domain)
	fmt.Println(internal.DeleteBodyWantJSON(path, body))

}
