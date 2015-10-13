package service

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdServiceShow = &cobra.Command{
	Use:     "show",
	Aliases: []string{"inspect"},
	Short:   "Show a docker service: sail service show [<applicationName>/]<serviceId>",
	Long: `Show a docker service: sail service show [<applicationName>/]<serviceId>
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
	// Split namespace and service
	host, app, service, _, err := internal.ParseResourceName(serviceID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	}

	internal.FormatOutputDef(internal.GetWantJSON(fmt.Sprintf("/applications/%s/services/%s", app, service)))
}
