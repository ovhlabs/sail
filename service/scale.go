package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/runabove/sail/internal"
)

var scaleBatch bool
var scaleDestroy bool
var scaleNumber int
var scaleUsage = "usage: sail services scale [-h] [--number NUMBER] [--batch] [--destroy] [<application>/]<service>"
var scaleLongUsage = `usage: sail services scale [-h] [--number NUMBER] [--batch] [--destroy] [<application>/]<service>

The command will exit as soon as all service containers have stopped.
Its exit status will be the one of the last container. If the last container was stopped with
a signal, the command exits with an exit status of 255.`

// Scale json data arguments
type Scale struct {
	Number  int  `json:"container_number"`
	Destroy bool `json:"destroy"`
}

func scaleCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "scale",
		Short: scaleUsage,
		Long:  scaleUsage,
		Run:   cmdScale,
	}

	cmd.Flags().BoolVar(&scaleBatch, "batch", false, "do not attach console on start")
	cmd.Flags().BoolVar(&scaleDestroy, "destroy", false, "when scaling down, prune last stopped containers")
	cmd.Flags().IntVar(&scaleNumber, "number", 0, "scale to `number` of containers")

	return cmd
}

func cmdScale(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, scaleUsage)
		os.Exit(1)
	}

	// Split namespace and service
	host, app, service, _, err := internal.ParseResourceName(args[0])
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	}

	serviceScale(app, service, scaleNumber, scaleDestroy, scaleBatch)
}

// serviceScale start service (without attach)
func serviceScale(app string, service string, number int, destroy bool, batch bool) {
	if !batch {
		internal.StreamPrint("GET", fmt.Sprintf("/applications/%s/services/%s/attach", app, service), nil)
	}

	// stream service events in a goroutine
	internal.EventStreamPrint("GET", fmt.Sprintf("/applications/%s/services/%s/events", app, service), nil, true)

	path := fmt.Sprintf("/applications/%s/services/%s/scale", app, service)

	args := Scale{
		Number:  number,
		Destroy: destroy,
	}

	data, err := json.Marshal(&args)
	internal.Check(err)

	buffer, _, err := internal.Stream("POST", path, data)
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
