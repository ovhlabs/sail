package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sail/internal"
)

var (
	logsBody Logs
)

func logsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Logs of a docker service : sail service logs <applicationName>/<serviceId>",
		Long:  `Logs of a docker service : sail service logs <applicationName>/<serviceId>`,
		Run:   cmdLogs,
	}

	cmd.Flags().IntVarP(&logsBody.Tail, "tail", "", 0, "Return N last lines, before offset")
	cmd.Flags().IntVarP(&logsBody.Head, "head", "", 0, "Return N first lines, after offset")
	cmd.Flags().IntVarP(&logsBody.Offset, "offset", "", 0, "Offset result by N line")
	cmd.Flags().StringVarP(&logsBody.Period, "period", "", "24 hours ago", "Lucene compatible period")
	cmd.Flags().StringVarP(&logsBody.Search, "search", "", "", "Only return matching lines")

	return cmd
}

// Logs struct holds all parameters sent to /applications/%s/services/%s/logs
type Logs struct {
	Application string `json:"-"`
	Service     string `json:"-"`

	Repository string `json:"repository,omitempty"`
	Tail       int    `json:"tail,omitempty"`
	Head       int    `json:"head,omitempty"`
	Offset     int    `json:"offset,omitempty"`
	Period     string `json:"period,omitempty"`
	Search     string `json:"search,omitempty"`
}

func cmdLogs(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail service logs <applicationName>/<serviceId>. Please see sail service logs --help\n"
	if len(args) != 1 {
		fmt.Printf(usage)
		return
	}

	split := strings.Split(args[0], "/")
	if len(split) != 2 {
		fmt.Printf(usage)
		return
	}

	// Get args
	logsBody.Application = split[0]
	logsBody.Service = split[1]
	serviceLogs(logsBody)
}

func serviceLogs(args Logs) {

	path := fmt.Sprintf("/applications/%s/services/%s/logs", args.Application, args.Service)
	body, err := json.MarshalIndent(args, " ", " ")
	if err != nil {
		fmt.Printf("Fatal: %s\n", err)
		return
	}

	fmt.Println(internal.ReqWantJSON("GET", http.StatusOK, path, body))
}
