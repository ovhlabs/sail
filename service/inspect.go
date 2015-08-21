package service

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sailgo/internal"
)

var cmdServiceInspect = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a docker service : sailgo service inspect <applicationName>/<serviceId>",
	Long: `Inspect a docker service : sailgo service inspect <applicationName>/<serviceId>
	\"example : sailgo service inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid usage. sailgo service inspect <applicationName>/<serviceId>. Please see sailgo service inspect --help")
		} else {
			serviceInspect(args[0])
		}
	},
}

func serviceInspect(serviceID string) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sailgo service inspect <applicationName>/<serviceId>. Please see sailgo service inspect --help")
	} else {
		fmt.Println(internal.GetWantJSON(fmt.Sprintf("/applications/%s/services/%s", t[0], t[1])))
	}
}
