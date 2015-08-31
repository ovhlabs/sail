package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sail/internal"
)

var startBatch bool
var startUsage = "usage: sail services start [-h] [--batch] service"

func startCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: startUsage,
		Long:  startUsage,
		Run:   cmdStart,
	}

	cmd.Flags().BoolVar(&startBatch, "batch", false, "do not attach console on start")

	return cmd
}

func cmdStart(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, startUsage)
		os.Exit(1)
	}

	t := strings.Split(args[0], "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, startUsage)
		os.Exit(1)
	}

	if batch {
		serviceStartRun(t[0], t[1])
		os.Exit(0)
		return
	}
	serviceStart(t[0], t[1], startBatch)
}

// serviceStart attach and start service
func serviceStart(app string, service string, batch bool) {

	reader, _, e := internal.Stream("GET",
		fmt.Sprintf("/applications/%s/services/%s/attach", app, service),
		nil,
		internal.SetHeader("Content-Type", "application/x-yaml"))

	if e != nil {
		internal.Exit("Error while attach: %s\n", e)
	}

	serviceStartRun(app, service)

	// Display api stream
	err := internal.DisplayStream(reader)
	if err != nil {
		internal.Exit("Error: %s\n", err)
	}
}

// serviceStart start service (without attach)
func serviceStartRun(app string, service string) {
	path := fmt.Sprintf("/applications/%s/services/%s/start?stream", app, service)
	buffer, _, err := internal.Stream("POST", path, []byte("{}"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	err = internal.DisplayStream(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
