package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sailgo/internal"
)

var cmdServiceAttach = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a service console : sailgo service attach <applicationName>/<serviceId>",
	Long: `Attach to a service console : sailgo service attach <applicationName>/<serviceId>
	\"example : sailgo service attach myApp myServiceId"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid usage. sailgo service attach <applicationName>/<serviceId>. Please see sailgo service attach --help")
		} else {
			serviceAttach(args[0])
		}
	},
}

func serviceAttach(serviceID string) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sailgo service attach <applicationName>/<serviceId>. Please see sailgo service attach --help")
	} else {
		internal.StreamWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/services/%s/attach", t[0], t[1]), nil)
	}
}
