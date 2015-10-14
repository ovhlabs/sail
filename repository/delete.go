package repository

import (
	"fmt"
	"os"

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
	// Split namespace and repository
	host, app, repo, tag, err := internal.ParseResourceName(repositoryID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	} else if len(tag) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid repository name. Please see sail repository delete --help\n")
		os.Exit(1)
	}

	path := fmt.Sprintf("/repositories/%s/%s", app, repo)
	data := internal.DeleteWantJSON(path)

	internal.FormatOutput(data, func(data []byte) {
		fmt.Fprintf(os.Stderr, "Deleted repository %s/%s\n", app, repo)
	})
}
