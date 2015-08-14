package main

import (
	"fmt"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"
)

func init() {
	cmdMe.AddCommand(cmdMeShow)
	cmdMe.AddCommand(cmdMeSetAcls)
}

var cmdMe = &cobra.Command{
	Use:     "me",
	Short:   "Me commands : sailgo me --help",
	Long:    `Me commands : sailgo me <command>`,
	Aliases: []string{"m"},
}

var cmdMeShow = &cobra.Command{
	Use:   "show",
	Short: "Display information about me : sailgo me show",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO sailgo me show
		fmt.Println("sailgo me show TO BE IMPLEMENTED")
	},
}

var cmdMeSetAcls = &cobra.Command{
	Use:   "setAcls",
	Short: "Set Acls : sailgo me setAcls",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO  sailgo me setAcls
		fmt.Println("sailgo me setAcls TO BE IMPLEMENTED")
	},
}
