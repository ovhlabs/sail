package compose

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	Cmd.AddCommand(cmdComposeUp)
	Cmd.AddCommand(cmdComposeGet)
}

// Cmd compose
var Cmd = &cobra.Command{
	Use:     "compose",
	Short:   "Compose commands : sailgo compose --help",
	Long:    `Compose commands : sailgo compose <command>`,
	Aliases: []string{"comp"},
}

var cmdComposeUp = &cobra.Command{
	Use:   "up",
	Short: "sailgo compose up <namespace>",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO sailgo compose up
		fmt.Println("sailgo compose up TO BE IMPLEMENTED")
	},
}

var cmdComposeGet = &cobra.Command{
	Use:   "get",
	Short: "sailgo compose get <namespace>",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO  sailgo compose get
		fmt.Println("sailgo compose get TO BE IMPLEMENTED")
	},
}
