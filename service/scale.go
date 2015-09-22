package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/runabove/sail/internal"
)

var scaleBatch bool
var scaleDestroy bool
var scaleNumber int
var scaleUsage = "usage: sail services scale [-h] [--number NUMBER] [--batch] [--destroy] <application>/<service>"

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

	t := strings.Split(args[0], "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, startUsage)
		os.Exit(1)
	}

	if scaleBatch {
		serviceScale(t[0], t[1], scaleNumber, scaleDestroy)
	} else {
		serviceScaleStream(t[0], t[1], scaleNumber, scaleDestroy)
	}
}

// serviceScaletStream attach and start service
func serviceScaleStream(app string, service string, number int, destroy bool) {

	reader, _, e := internal.Stream("GET",
		fmt.Sprintf("/applications/%s/services/%s/attach", app, service),
		nil,
		internal.SetHeader("Content-Type", "application/x-yaml"))

	if e != nil {
		internal.Exit("Error while attach: %s\n", e)
	}

	serviceScale(app, service, number, destroy)

	// Display api stream
	err := internal.DisplayStream(reader)
	internal.Check(err)
}

// serviceScale start service (without attach)
func serviceScale(app string, service string, number int, destroy bool) {
	path := fmt.Sprintf("/applications/%s/services/%s/scale?stream", app, service)

	args := Scale{
		Number:  number,
		Destroy: destroy,
	}

	data, err := json.Marshal(&args)
	internal.Check(err)

	buffer, _, err := internal.Stream("POST", path, data)
	internal.Check(err)

	err = internal.DisplayStream(buffer)
	internal.Check(err)
}
