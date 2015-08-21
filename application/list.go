package application

import (
	"fmt"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sailgo/internal"
)

var cmdApplicationList = &cobra.Command{
	Use:     "list",
	Short:   "List granted apps : sailgo application list",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(internal.GetWantJSON("/applications"))
	},
}
