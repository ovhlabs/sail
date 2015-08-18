package version

import (
	"fmt"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sailgo/internal"
)

var Cmd = &cobra.Command{
	Use:     "version",
	Short:   "Display Version of sailgo : sailgo version",
	Long:    `sailgo version`,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version sailgo : %s\n", internal.VERSION)
		internal.ReadConfig()
	},
}
