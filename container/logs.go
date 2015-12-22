package container

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/google/go-querystring/query"
	"github.com/spf13/cobra"

	"github.com/runabove/sail/internal"
)

var (
	logsBody Logs
)

func cmdContainerLogs() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Fetch the logs of a container",
		Long:  "Fetch the logs of a container",
		Run:   cmdLogs,
	}

	cmd.Flags().IntVarP(&logsBody.Tail, "tail", "", 0, "Return N last lines, before offset.")
	cmd.Flags().IntVarP(&logsBody.Head, "head", "", 0, "Return N first lines, after offset.")
	cmd.Flags().IntVarP(&logsBody.Offset, "offset", "", 0, "Offset result by N line")
	cmd.Flags().StringVarP(&logsBody.Period, "period", "", "", "Human readable (Lucene syntax) period")

	return cmd
}

// Logs struct holds all parameters sent to /applications/%s/containers/%s/logs
type Logs struct {
	Application string `url:"-"`
	Container   string `url:"-"`

	Tail   int    `url:"tail,omitempty"`
	Head   int    `url:"head,omitempty"`
	Offset int    `url:"offset,omitempty"`
	Period string `url:"period,omitempty"`
}

func cmdLogs(cmd *cobra.Command, args []string) {
	usage := "usage: sail containers logs [<applicationName>/]<containerId>"

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	// Split namespace and container
	host, app, container, tag, err := internal.ParseResourceName(args[0])
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	} else if len(tag) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid container name. Please see sail container logs --help\n")
		os.Exit(1)
	}

	// Get args
	logsBody.Application = app
	logsBody.Container = container
	containerLogs(logsBody)
}

func containerLogs(args Logs) {
	queryArgs, err := query.Values(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err)
		return
	}
	path := fmt.Sprintf("/applications/%s/containers/%s/logs?%s", args.Application, args.Container, queryArgs.Encode())

	b := internal.GetWantJSON(path)
	internal.FormatOutput(b, containerLogsFormatter)
}

func containerLogsFormatter(data []byte) {
	var logs [][]string
	err := json.Unmarshal(data, &logs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	titles := []string{"TIMESTAMP", "ID", "LOG"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	for i := range logs {
		fmt.Fprintf(w, "%s\t%s\t%s\n", logs[i][0], logs[i][1], logs[i][2])
		w.Flush()
	}

}
