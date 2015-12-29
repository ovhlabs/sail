package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/runabove/sail/internal"
)

var startBatch bool
var startUsage = "usage: sail services start [-h] [--batch] [<applicationName>/]<serviceId>"

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

	// Split namespace and service
	host, app, service, _, err := internal.ParseResourceName(args[0])
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	}

	serviceStart(app, service, startBatch)
}

// serviceStart start service (without attach)
func serviceStart(app string, service string, batch bool) {
	if !batch {
		internal.StreamPrint("GET", fmt.Sprintf("/applications/%s/services/%s/attach", app, service), nil)
	}

	path := fmt.Sprintf("/applications/%s/services/%s/start?stream", app, service)
	buffer, _, err := internal.Stream("POST", path, []byte("{}"))
	internal.Check(err)

	line, err := internal.DisplayStream(buffer)
	internal.Check(err)
	if len(line) > 0 {
		var data map[string]interface{}
		err = json.Unmarshal(line, &data)
		internal.Check(err)

		fmt.Printf("Hostname: %v\n", data["hostname"])
	}

	if !batch {
		internal.ExitAfterCtrlC()
	}
}
