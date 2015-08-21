package network

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdNetworkList)
	Cmd.AddCommand(cmdNetworkInspect)

	// TODO
	// sail networks add        Add a new private network
	// sail networks range-add  Add an allocation range to a private network
	// sail networks rm         Delete a private network

}

// Cmd network
var Cmd = &cobra.Command{
	Use:     "network",
	Short:   "Network commands : sailgo network --help",
	Long:    `Network commands : sailgo network <command>`,
	Aliases: []string{"networks"},
}
