package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmdMe.AddCommand(cmdMeShow)
	cmdMe.AddCommand(cmdMeSetAcls)
}

var cmdMe = &cobra.Command{
	Use:     "group",
	Short:   "Me commands : sailgo me --help",
	Long:    `Me commands : sailgo me <command>`,
	Aliases: []string{"m"},
}

var cmdMeShow = &cobra.Command{
	Use:   "show",
	Short: "Display information about me : sailgo me show",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		fmt.Println("TO BE IMPLEMENTED")
	},
}

var cmdMeSetAcls = &cobra.Command{
	Use:   "show",
	Short: "Set Acls : sailgo me setAcls",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
		fmt.Println("TO BE IMPLEMENTED")
	},
}
