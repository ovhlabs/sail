package me

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdMeShow)
	Cmd.AddCommand(cmdMeSetAcl)
}

// Cmd me
var Cmd = &cobra.Command{
	Use:     "me",
	Short:   "Me commands : sailgo me --help",
	Long:    `Me commands : sailgo me <command>`,
	Aliases: []string{"m"},
}
