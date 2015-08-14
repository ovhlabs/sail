package main

import "stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"

func init() {
	// TODO sailgo repository
}

var cmdRepository = &cobra.Command{
	Use:     "repository",
	Short:   "Repository commands : sailgo repository --help",
	Long:    `Repository commands : sailgo repository <command>`,
	Aliases: []string{"r", "repo", "repositories"},
}
