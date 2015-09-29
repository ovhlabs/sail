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

	serviceScale(t[0], t[1], scaleNumber, scaleDestroy, scaleBatch)
}

// serviceScale start service (without attach)
func serviceScale(app string, service string, number int, destroy bool, batch bool) {
	if !batch {
		internal.StreamPrint("GET", fmt.Sprintf("/applications/%s/services/%s/attach", app, service), nil)
	}

	path := fmt.Sprintf("/applications/%s/services/%s/scale?stream", app, service)

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
	if line != nil {
		var data map[string]interface{}
		err = json.Unmarshal(line, &data)
		internal.Check(err)

		fmt.Printf("Hostname: %v\n", data["hostname"])
		fmt.Printf("Running containers: %v/%v\n", data["container_number"], data["container_target"])
	}

	if !batch {
		internal.ExitAfterCtrlC()
	}
}
