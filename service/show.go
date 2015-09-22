package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/runabove/sail/internal"
)

var cmdServiceShow = &cobra.Command{
	Use:     "show",
	Aliases: []string{"inspect"},
	Short:   "Show a docker service: sail service show <applicationName>/<serviceId>",
	Long: `Show a docker service: sail service show <applicationName>/<serviceId>
	\"example: sail service show my-app"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail service show <applicationName>/<serviceId>. Please see sail service show --help")
		} else {
			serviceShow(args[0])
		}
	},
}

func serviceShow(serviceID string) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid usage. sail service show <applicationName>/<serviceId>. Please see sail service show --help")
	} else {
		internal.FormatOutputDef(internal.GetWantJSON(fmt.Sprintf("/applications/%s/services/%s", t[0], t[1])))
	}
}
