package operation

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdOperationList = &cobra.Command{
	Use:     "list",
	Short:   "List the services: sail operation list [applicationName]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		operationList(internal.GetListApplications(args))
	},
}

func operationList(apps []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	titles := []string{"APPLICATION", "SERVICE", "ID", "COMMAND", "SUBMITTED (UTC)"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	operations := []map[string]string{}
	for _, app := range apps {
		// Sanity checks
		err := internal.CheckName(app)
		internal.Check(err)

		r := internal.GetWantJSON(fmt.Sprintf("/operation/application/%s", app))
		internal.Check(json.Unmarshal(r, &operations))
		for _, operation := range operations {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				app,
				operation["service"],
				operation["topic"],
				operation["command"],
				operation["started_at"])
		}
	}
	w.Flush()
}
