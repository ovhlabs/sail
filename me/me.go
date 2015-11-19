package me

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdMeShow)
	Cmd.AddCommand(cmdMeSetPassword)
	Cmd.AddCommand(cmdMeSetAcl)

	cmdMeSSHKey.AddCommand(cmdMeSSHKeyList)
	cmdMeSSHKey.AddCommand(cmdMeSSHKeyAdd)
	cmdMeSSHKey.AddCommand(cmdMeSSHKeyDelete)

	Cmd.AddCommand(cmdMeSSHKey)
}

// Cmd me
var Cmd = &cobra.Command{
	Use:     "me",
	Short:   "Me commands: sail me --help",
	Long:    `Me commands: sail me <command>`,
	Aliases: []string{"m"},
}
