package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/runabove/sail/internal"
)

var cmdServiceStop = &cobra.Command{
	Use:   "stop",
	Short: "Stop a docker service: sail service stop <applicationName>/<serviceId>",
	Long: `Stop a docker service: sail service stop <applicationName>/<serviceId>
	\"example: sail service stop my-app/myService"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail service stop <applicationName>/<serviceId>. Please see sail service stop --help")
		} else {
			serviceStop(args[0])
		}
	},
}

func serviceStop(serviceID string) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid usage. sail service stop <applicationName>/<serviceId>. Please see sail service stop --help")
	} else {
		var empty map[string]interface{}
		em, _ := json.Marshal(empty)
		//TODO: attach + print stream (stop logs)
		internal.FormatOutputDef(internal.PostBodyWantJSON(fmt.Sprintf("/applications/%s/services/%s/stop", t[0], t[1]), em))
	}
}
