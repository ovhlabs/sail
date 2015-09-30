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
		Short: "sail compose get <application>",
		Run:   cmdGet,
	}

	cmd.Flags().BoolVarP(&getStandard, "standard", "", false, "Return only Docker Compose standard properties")
	return cmd
}

func cmdGet(cmd *cobra.Command, args []string) {
	// FIXME: duplicate
	internal.ReadConfig()
	var ns string

	// Check args
	if len(args) > 1 {
		internal.Exit("Invalid usage. sail compose get [--standard] [<application>]. Please see sail compose get -h\n")
	} else if len(args) > 1 {
		ns = args[0]
	} else {
		ns = internal.User
	}

	path := fmt.Sprintf("/applications/%s/fig?standard=%v", ns, getStandard)
	data, _, err := internal.Request("GET", path, nil)
	if err != nil {
		internal.Exit("Error: %s\n", err)
	}

	fmt.Printf("%s", data)
}
