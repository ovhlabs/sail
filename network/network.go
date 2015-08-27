package network

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdNetworkAdd)
	Cmd.AddCommand(cmdNetworkShow)
	Cmd.AddCommand(cmdNetworkList)
	Cmd.AddCommand(cmdNetworkRangeAdd)
	Cmd.AddCommand(cmdNetworkDelete)
}

// Cmd network
var Cmd = &cobra.Command{
	Use:     "network",
	Short:   "Network commands: sail network --help",
	Long:    `Network commands: sail network <command>`,
	Aliases: []string{"networks", "net"},
}
