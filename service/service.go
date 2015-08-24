package service

import (
	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sailgo/service/domain"
)

func init() {
	Cmd.AddCommand(cmdServiceAttach)
	Cmd.AddCommand(cmdServiceList)
	Cmd.AddCommand(cmdServiceInspect)
	Cmd.AddCommand(cmdServiceStop)
	Cmd.AddCommand(domain.Cmd)
	Cmd.AddCommand(logsCmd())
	Cmd.AddCommand(redeployCmd())
	Cmd.AddCommand(addCmd())
	Cmd.AddCommand(rmCmd())
	Cmd.AddCommand(startCmd())
	Cmd.AddCommand(scaleCmd())
}

// Cmd service
var Cmd = &cobra.Command{
	Use:     "service",
	Short:   "Service commands : sailgo service --help",
	Long:    `Service commands : sailgo service <command>`,
	Aliases: []string{"services"},
}
