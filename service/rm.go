package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sail/internal"
)

func rmCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "rm",
		Short:   "Remove a docker service: sail service rm <applicationName>/<serviceId>",
		Run:     cmdServiceRm,
		Aliases: []string{"delete", "remove"},
	}
	return cmd
}

func cmdServiceRm(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail service rm <applicationName>/<serviceId>"
	if len(args) == 0 {
		fmt.Println(usage)
		return
	}
	argsS := strings.Split(args[0], "/")
	if len(argsS) != 2 {
		fmt.Println(usage)
		return
	}

	serviceRm(argsS[0], argsS[1])
}

func serviceRm(namespace string, name string) {
	path := fmt.Sprintf("/applications/%s/services/%s", namespace, name)
	data := internal.ReqWant("DELETE", http.StatusOK, path, nil)
	// TODO Check for json error here
	fmt.Printf("%s\n", data)
}
