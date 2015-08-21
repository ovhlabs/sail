package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sailgo/internal"
)

var versionNewLine bool

func init() {
	Cmd.Flags().BoolVarP(&versionNewLine, "versionNewLine", "", true, "New line after version number.")
}

// Cmd version
var Cmd = &cobra.Command{
	Use:     "version",
	Short:   "Display Version of sailgo : sailgo version",
	Long:    `sailgo version`,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		if versionNewLine {
			fmt.Printf("Version sailgo : %s\n", internal.VERSION)
			internal.ReadConfig()
		} else {
			fmt.Printf(internal.VERSION)
		}
	},
}
