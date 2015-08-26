package container

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sail/internal"
)

var logsTail int
var logsHead int
var logsTimestamp bool

func cmdContainerLogs() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Fetch the logs of a container",
		Long:  "Fetch the logs of a container",
		Run:   cmdLogs,
	}

	// Apparently not in the latest version of sail
	//cmd.Flags().IntVar(&logsTail, "tail", 10, "Return N last lines, before offset.")
	//cmd.Flags().IntVar(&logsHead, "head", 10, "Return N first lines, after offset.")
	//cmd.Flags().BoolVarP(&logsTimestamp, "timestamp", "t", false, "filter offset.")

	return cmd
}

func cmdLogs(cmd *cobra.Command, args []string) {
	usage := "usage: sail containers logs <applicationName>/<containerId>"

	if len(args) != 1 {
		fmt.Println(usage)
		os.Exit(1)
	}

	t := strings.Split(args[0], "/")
	if len(t) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	containerLogs(t[0], t[1], logsHead, logsTail, logsTimestamp)
}

func containerLogs(app string, container string, head int, tail int, ts bool) {
	path := fmt.Sprintf("/applications/%s/containers/%s/logs?timestamps=%v", app, container, ts)

	data, _, err := internal.Request("GET", path, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	var logs [][]string
	err = json.Unmarshal(data, &logs)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	titles := []string{"TIMESTAMP", "ID", "LOG"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	for i := range logs {
		fmt.Fprintf(w, "%s\t%s\t%s\n", logs[i][0], logs[i][1], logs[i][2])
		w.Flush()
	}
}
