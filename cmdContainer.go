package main

import "stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"

func init() {
	//TODO
}

var cmdContainer = &cobra.Command{
	Use:     "container",
	Short:   "Container commands : sailgo container --help",
	Long:    `Container commands : sailgo container <command>`,
	Aliases: []string{"c", "containers"},
}
