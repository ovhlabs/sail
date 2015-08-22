package network

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdNetworkAdd)
	Cmd.AddCommand(cmdNetworkInspect)
	Cmd.AddCommand(cmdNetworkList)
	Cmd.AddCommand(cmdNetworkRangeAdd)
	Cmd.AddCommand(cmdNetworkRm)
}

// Cmd network
var Cmd = &cobra.Command{
	Use:     "network",
	Short:   "Network commands : sailgo network --help",
	Long:    `Network commands : sailgo network <command>`,
	Aliases: []string{"networks", "net"},
}
