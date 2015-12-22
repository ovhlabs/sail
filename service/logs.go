package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/google/go-querystring/query"
	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var (
	logsBody Logs
)

func logsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Logs of a docker service: sail service logs [<applicationName>/]<serviceId>",
		Long:  `Logs of a docker service: sail service logs [<applicationName>/]<serviceId>`,
		Run:   cmdLogs,
	}

	cmd.Flags().IntVarP(&logsBody.Tail, "tail", "", 0, "Return N last lines, before offset")
	cmd.Flags().IntVarP(&logsBody.Head, "head", "", 0, "Return N first lines, after offset")
	cmd.Flags().IntVarP(&logsBody.Offset, "offset", "", 0, "Offset result by N line")
	cmd.Flags().StringVarP(&logsBody.Period, "period", "", "", "Human readable (Lucene syntax) period")
	cmd.Flags().StringVarP(&logsBody.Search, "search", "", "", "Only return matching lines")

	return cmd
}

// Logs struct holds all parameters sent to /applications/%s/services/%s/logs
type Logs struct {
	Application string `url:"-"`
	Service     string `url:"-"`

	Repository string `url:"repository,omitempty"`
	Tail       int    `url:"tail,omitempty"`
	Head       int    `url:"head,omitempty"`
	Offset     int    `url:"offset,omitempty"`
	Period     string `url:"period,omitempty"`
	Search     string `url:"search,omitempty"`
}

func cmdLogs(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail service logs [<applicationName>/]<serviceId>. Please see sail service logs --help\n"
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	// Split namespace and service
	host, app, service, _, err := internal.ParseResourceName(args[0])
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	}

	// Get args
	logsBody.Application = app
	logsBody.Service = service
	serviceLogs(logsBody)
}

func serviceLogs(args Logs) {

	queryArgs, err := query.Values(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err)
		return
	}
	path := fmt.Sprintf("/applications/%s/services/%s/logs?%s", args.Application, args.Service, queryArgs.Encode())

	b := internal.GetWantJSON(path)
	internal.FormatOutput(b, serviceLogsFormatter)
}

func serviceLogsFormatter(data []byte) {
	logs := [][]string{}
	err := json.Unmarshal(data, &logs)
	internal.Check(err)

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	titles := []string{"TIMESTAMP", "ID", "LOG"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	for i := range logs {
		fmt.Fprintf(w, "%s\t%s\t%s\n", logs[i][0], logs[i][1], logs[i][2])
		w.Flush()
	}
}
