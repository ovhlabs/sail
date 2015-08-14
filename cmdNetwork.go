package main

import "stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"

func init() {
	// TODO sail go network
}

var cmdNetwork = &cobra.Command{
	Use:     "network",
	Short:   "Network commands : sailgo network --help",
	Long:    `Network commands : sailgo network <command>`,
	Aliases: []string{"networks"},
}
