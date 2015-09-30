package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/runabove/sail/internal"
)

func deleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a docker service: sail service delete [<applicationName>/]<serviceId>",
		Run:     cmdServiceDelete,
		Aliases: []string{"del", "rm", "remove"},
	}
	return cmd
}

func cmdServiceDelete(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail service delete [<applicationName>/]<serviceId>"
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, usage)
		return
	}

	// Split namespace and service
	host, app, service, _, err := internal.ParseResourceName(args[0])
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	}

	serviceDelete(app, service)
}

func serviceDelete(namespace string, name string) {
	path := fmt.Sprintf("/applications/%s/services/%s", namespace, name)
	data := internal.ReqWant("DELETE", http.StatusOK, path, nil)
	// TODO Check for json error here
	internal.FormatOutputDef(data)
}
