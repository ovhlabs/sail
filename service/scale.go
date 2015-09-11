package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sail/internal"
)

var scaleBatch bool
var scaleNumber int
var scaleUsage = "usage: sail services scale [-h] [--number NUMBER] [--batch] <application>/<service>"

// Scale json data arguments
type Scale struct {
	Number int `json:"container_number"`
}

func scaleCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "scale",
		Short: scaleUsage,
		Long:  scaleUsage,
		Run:   cmdScale,
	}

	cmd.Flags().BoolVar(&scaleBatch, "batch", false, "do not attach console on start")
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

	serviceScale(t[0], t[1], scaleNumber, scaleBatch)
}

func serviceScale(app string, service string, number int, batch bool) {
	path := fmt.Sprintf("/applications/%s/services/%s/scale?stream", app, service)

	args := Scale{
		Number: number,
	}

	data, err := json.Marshal(&args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	buffer, _, err := internal.Stream("POST", path, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	err = internal.DisplayStream(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	serviceStart(app, service, batch)
}
