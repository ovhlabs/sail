package repository

import (
	"fmt"
	"os"
	"strings"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

var cmdRepositoryDelete = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm", "remove"},
	Short:   "Delete a repository: sail repository delete <applicationName>/<repositoryId>",
	Long:    `Delete a repository: sail repository delete <applicationName>/<repositoryId>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail repository delete <applicationName>/<repositoryId>. Please see sail repository delete --help")
		} else {
			repositoryRemove(args[0])
		}
	},
}

func repositoryRemove(repositoryID string) {
	t := strings.Split(repositoryID, "/")
	if len(t) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid usage. sail repository delete <applicationName>/<repositoryId>. Please see sail repository delete --help")
		return
	}

	path := fmt.Sprintf("/repositories/%s/%s", t[0], t[1])
	internal.FormatOutputDef(internal.DeleteWantJSON(path))
}
