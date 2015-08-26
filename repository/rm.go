package repository

import (
	"fmt"
	"os"
	"strings"

	"stash.ovh.net/sailabove/sail/internal"

	"github.com/spf13/cobra"
)

var cmdRepositoryRm = &cobra.Command{
	Use:     "rm",
	Short:   "Delete a repository : sail repository rm <applicationName>/<repositoryId>",
	Long:    `Delete a repository : sail repository rm <applicationName>/<repositoryId>`,
	Aliases: []string{"delete", "remove"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail repository rm <applicationName>/<repositoryId>. Please see sail repository rm --help")
		} else {
			repositoryRemove(args[0])
		}
	},
}

func repositoryRemove(repositoryID string) {
	t := strings.Split(repositoryID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid usage. sail repository rm <applicationName>/<repositoryId>. Please see sail repository rm --help")
		return
	}

	path := fmt.Sprintf("/repositories/%s/%s", t[0], t[1])
	internal.FormatOutputDef(internal.DeleteWantJSON(path))
}
