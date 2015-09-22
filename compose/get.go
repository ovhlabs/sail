package compose

import (
	"fmt"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

var getStandard bool

func cmdComposeGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "sail compose get <namespace>",
		Run:   cmdGet,
	}

	cmd.Flags().BoolVarP(&getStandard, "standard", "", false, "Return only Docker Compose standard properties")
	return cmd
}

func cmdGet(cmd *cobra.Command, args []string) {

	// Check arguments
	if len(args) != 1 {
		internal.Exit("Invalid usage. sail compose get [-h] [--standard] namespace\n")
	}

	path := fmt.Sprintf("/applications/%s/fig?standard=%v", args[0], getStandard)
	data, _, err := internal.Request("GET", path, nil)
	if err != nil {
		internal.Exit("Error: %s\n", err)
	}

	fmt.Printf("%s", data)
}
