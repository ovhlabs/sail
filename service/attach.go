package service

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdServiceAttach = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a service console: sail service attach <applicationName>/<serviceId>",
	Long: `Attach to a service console: sail service attach <applicationName>/<serviceId>
	\"example: sail service attach my-app myServiceId"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail service attach <applicationName>/<serviceId>. Please see sail service attach --help")
		} else {
			serviceAttach(args[0])
		}
	},
}

func serviceAttach(serviceID string) {
	// Split namespace and service
	host, app, service, _, err := internal.ParseResourceName(serviceID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	}

	internal.StreamPrint("GET", fmt.Sprintf("/applications/%s/services/%s/attach", app, service), nil)
	internal.ExitAfterCtrlC()
}
