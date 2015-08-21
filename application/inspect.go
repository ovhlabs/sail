package application

import (
	"fmt"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sailgo/internal"
)

var cmdApplicationInspect = &cobra.Command{
	Use:   "inspect",
	Short: "Details of an app : sailgo application inspect <applicationName>",
	Long: `Details of an app : sailgo application inspect <applicationName>
	\"example : sailgo application inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Println("Invalid usage. Please see sailgo application inspect --help")
		} else {
			fmt.Println(internal.GetWantJSON(fmt.Sprintf("/applications/%s", args[0])))
		}
	},
}
