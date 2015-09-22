package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/runabove/sail/internal"
)

var cmdServiceAttach = &cobra.Command{
	Use:   "attach",
	Short: "Attach to a service console: sail service attach <applicationName>/<serviceId>",
	Long: `Attach to a service console: sail service attach <applicationName>/<serviceId>
	\"example: sail service attach my-app myServiceId"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail service attach <applicationName>/<serviceId>. Please see sail service attach --help")
		} else {
			serviceAttach(args[0])
		}
	},
}

func serviceAttach(serviceID string) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid usage. sail service attach <applicationName>/<serviceId>. Please see sail service attach --help")
	} else {
		internal.StreamWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/services/%s/attach", t[0], t[1]), nil)
	}
}
