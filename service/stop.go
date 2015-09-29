package service

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var stopBatch bool
var stopUsage = "usage: sail services stop [-h] [--batch] [<applicationName>/]<serviceId>"

func stopCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "stop",
		Short: stopUsage,
		Long:  stopUsage,
		Run:   cmdStop,
	}

	cmd.Flags().BoolVar(&startBatch, "batch", false, "do not attach console on stop")

	return cmd
}

func cmdStop(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, stopUsage)
		os.Exit(1)
	}

	// Split namespace and service
	host, app, service, _, err := internal.ParseResourceName(args[0])
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	}

	serviceStop(app, service, stopBatch)
}

// serviceStop stop service (without attach)
func serviceStop(app string, service string, batch bool) {
	if !batch {
		internal.StreamPrint("GET", fmt.Sprintf("/applications/%s/services/%s/attach", app, service), nil)
	}

	path := fmt.Sprintf("/applications/%s/services/%s/stop?stream", app, service)
	buffer, _, err := internal.Stream("POST", path, []byte("{}"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	_, err = internal.DisplayStream(buffer)
	internal.Check(err)
}
