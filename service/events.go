package service

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdServiceEvents = &cobra.Command{
	Use:   "events",
	Short: "Stream all service events: sail service events <applicationName>/<serviceId>",
	Long: `Stream all service events: sail service events <applicationName>/<serviceId>
	\"example: sail service events my-app/myServiceId"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail service events <applicationName>/<serviceId>. Please see sail service events --help")
		} else {
			serviceEvents(args[0])
		}
	},
}

func serviceEvents(serviceID string) {
	// Split namespace and service
	host, app, service, _, err := internal.ParseResourceName(serviceID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	}

	internal.EventStreamPrint("GET", fmt.Sprintf("/applications/%s/services/%s/events", app, service), nil, false)
	internal.ExitAfterCtrlC()
}
