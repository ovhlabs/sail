package service

import (
	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sail/service/domain"
)

func init() {
	Cmd.AddCommand(cmdServiceAttach)
	Cmd.AddCommand(cmdServiceList)
	Cmd.AddCommand(cmdServiceShow)
	Cmd.AddCommand(cmdServiceStop)
	Cmd.AddCommand(domain.Cmd)
	Cmd.AddCommand(logsCmd())
	Cmd.AddCommand(redeployCmd())
	Cmd.AddCommand(addCmd())
	Cmd.AddCommand(deleteCmd())
	Cmd.AddCommand(startCmd())
	Cmd.AddCommand(scaleCmd())
}

// Cmd service
var Cmd = &cobra.Command{
	Use:     "service",
	Short:   "Service commands: sail service --help",
	Long:    `Service commands: sail service <command>`,
	Aliases: []string{"services"},
}
