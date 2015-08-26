package me

import (
	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sail/internal"
)

var cmdMeShow = &cobra.Command{
	Use:   "show",
	Short: "Show account details : sail me show",
	Run: func(cmd *cobra.Command, args []string) {
		internal.FormatOutputDef(internal.GetWantJSON("/users"))
	},
}
