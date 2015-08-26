package domain

import (
	"encoding/json"
	"fmt"
	"strings"

	"stash.ovh.net/sailabove/sail/internal"

	"github.com/spf13/cobra"
)

var usageDomainAttach = "Invalid usage. sail service attach <applicationName>/<serviceId> <domain> <pattern> <method>. Please see sail service attach --help"
var cmdDomainAttach = &cobra.Command{
	Use:     "attach",
	Short:   "Attach a domain on the HTTP load balancer : sail service domain attach <applicationName>/<serviceId> <domain> <pattern> <method>",
	Aliases: []string{"add"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println(usageDomainAttach)
		} else {
			serviceDomainAttach(args[0], args[1], domainStruct{Pattern: args[2], Method: args[3]})
		}
	},
}

func serviceDomainAttach(serviceID, domain string, args domainStruct) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Println(usageDomainAttach)
		return
	}

	body, err := json.Marshal(args)
	internal.Check(err)

	path := fmt.Sprintf("/applications/%s/services/%s/attached-routes/%s", t[0], t[1], domain)
	fmt.Println(internal.PostBodyWantJSON(path, body))

}
