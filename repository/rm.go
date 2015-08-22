package repository

import (
	"fmt"
	"strings"

	"stash.ovh.net/sailabove/sailgo/internal"

	"github.com/spf13/cobra"
)

var cmdRepositoryRm = &cobra.Command{
	Use:     "rm",
	Short:   "Delete a repository : sailgo repository rm <applicationName>/<repositoryId>",
	Long:    `Delete a repository : sailgo repository rm <applicationName>/<repositoryId>`,
	Aliases: []string{"delete", "remove"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid usage. sailgo repository rm <applicationName>/<repositoryId>. Please see sailgo repository rm --help")
		} else {
			repositoryRemove(args[0])
		}
	},
}

func repositoryRemove(repositoryID string) {
	t := strings.Split(repositoryID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sailgo repository rm <applicationName>/<repositoryId>. Please see sailgo repository rm --help")
		return
	}

	path := fmt.Sprintf("/repositories/%s/%s", t[0], t[1])
	fmt.Println(internal.DeleteWantJSON(path))
}
