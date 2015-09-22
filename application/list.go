package application

import (
	"github.com/spf13/cobra"
	"github.com/runabove/sail/internal"
)

var cmdApplicationList = &cobra.Command{
	Use:     "list",
	Short:   "List granted apps: sail application list",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		internal.FormatOutputDef(internal.GetWantJSON("/applications"))
	},
}
