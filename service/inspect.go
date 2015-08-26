package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sail/internal"
)

var cmdServiceInspect = &cobra.Command{
	Use:     "inspect",
	Aliases: []string{"show"},
	Short:   "Inspect a docker service: sail service inspect <applicationName>/<serviceId>",
	Long: `Inspect a docker service: sail service inspect <applicationName>/<serviceId>
	\"example: sail service inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail service inspect <applicationName>/<serviceId>. Please see sail service inspect --help")
		} else {
			serviceInspect(args[0])
		}
	},
}

func serviceInspect(serviceID string) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid usage. sail service inspect <applicationName>/<serviceId>. Please see sail service inspect --help")
	} else {
		internal.FormatOutputDef(internal.GetWantJSON(fmt.Sprintf("/applications/%s/services/%s", t[0], t[1])))
	}
}
