package application

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdApplicationShow = &cobra.Command{
	Use:     "show",
	Aliases: []string{"inspect"},
	Short:   "Details of an app: sail application show <applicationName>",
	Long: `Details of an app: sail application show <applicationName>
	\"example: sail application show my-app"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application show --help")
		} else {
			// Sanity
			err := internal.CheckName(args[0])
			internal.Check(err)

			internal.FormatOutputDef(internal.GetWantJSON(fmt.Sprintf("/applications/%s", args[0])))
		}
	},
}
