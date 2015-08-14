package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VERSION of sailgo
const VERSION = "0.0.1"

var cmdVersion = &cobra.Command{
	Use:     "version",
	Short:   "Display Version of sailgo : sailgo version",
	Long:    `sailgo version`,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version sailgo : %s\n", VERSION)
		readConfig()
	},
}
