package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sail/internal"
)

var versionNewLine bool

func init() {
	Cmd.Flags().BoolVarP(&versionNewLine, "versionNewLine", "", true, "New line after version number.")
}

// Cmd version
var Cmd = &cobra.Command{
	Use:     "version",
	Short:   "Display Version of sail : sail version",
	Long:    `sail version`,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		if versionNewLine {
			fmt.Printf("Version sail : %s\n", internal.VERSION)
			internal.ReadConfig()
		} else {
			fmt.Printf(internal.VERSION)
		}
	},
}
