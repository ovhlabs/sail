package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/runabove/sail/internal"
)

var deleteForce bool

func deleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a docker service: sail service delete [<applicationName>/]<serviceId> [--force]",
		Run:     cmdServiceDelete,
		Aliases: []string{"del", "rm", "remove"},
	}

	cmd.Flags().BoolVarP(&deleteForce, "force", "", false, "danger zone: delete service even if it breaks links")

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
	path := fmt.Sprintf("/applications/%s/services/%s?force=%t", namespace, name, deleteForce)
	data := internal.ReqWant("DELETE", http.StatusOK, path, nil)

	internal.FormatOutput(data, func(data []byte) {
		fmt.Fprintf(os.Stderr, "Deleted service %s/%s\n", namespace, name)
	})
}
