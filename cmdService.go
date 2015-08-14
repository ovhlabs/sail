package main

import "stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"

func init() {
	// TODO sailgo service
}

var cmdService = &cobra.Command{
	Use:     "service",
	Short:   "Service commands : sailgo service --help",
	Long:    `Service commands : sailgo service <command>`,
	Aliases: []string{"s", "services"},
}
