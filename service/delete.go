package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sail/internal"
)

func deleteCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a docker service: sail service delete <applicationName>/<serviceId>",
		Run:     cmdServiceDelete,
		Aliases: []string{"del", "rm", "remove"},
	}
	return cmd
}

func cmdServiceDelete(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail service delete <applicationName>/<serviceId>"
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, usage)
		return
	}
	argsS := strings.Split(args[0], "/")
	if len(argsS) != 2 {
		fmt.Fprintln(os.Stderr, usage)
		return
	}

	serviceDelete(argsS[0], argsS[1])
}

func serviceDelete(namespace string, name string) {
	path := fmt.Sprintf("/applications/%s/services/%s", namespace, name)
	data := internal.ReqWant("DELETE", http.StatusOK, path, nil)
	// TODO Check for json error here
	internal.FormatOutputDef(data)
}
