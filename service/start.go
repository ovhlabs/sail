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
		fmt.Println(startUsage)
		os.Exit(1)
	}

	t := strings.Split(args[0], "/")
	if len(t) != 2 {
		fmt.Println(startUsage)
		os.Exit(1)
	}

	serviceStart(t[0], t[1], startBatch)
}

func serviceStart(app string, service string, batch bool) {
	path := fmt.Sprintf("/applications/%s/services/%s/start?stream", app, service)

	buffer, _, err := internal.Stream("POST", path, []byte("{}"))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	err = internal.DisplayStream(buffer)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if batch {
		os.Exit(0)
	}

	serviceAttach(app + "/" + service)
}
