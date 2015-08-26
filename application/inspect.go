package application

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sail/internal"
)

var cmdApplicationInspect = &cobra.Command{
	Use:     "inspect",
	Aliases: []string{"show"},
	Short:   "Details of an app: sail application inspect <applicationName>",
	Long: `Details of an app: sail application inspect <applicationName>
	\"example: sail application inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application inspect --help")
		} else {
			internal.FormatOutputDef(internal.GetWantJSON(fmt.Sprintf("/applications/%s", args[0])))
		}
	},
}
