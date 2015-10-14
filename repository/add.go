package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

var cmdRepositoryAdd = &cobra.Command{
	// FIXME: only support adding from source via CLI, type external. Rename to 'register' ?
	Use:   "add",
	Short: "Add a new repository: sail repository add [<applicationName>/]<repositoryId> <type> [source]",
	Long: `Add a new repository: sail repository add [<applicationName>/]<repositoryId> <type> [source]

	<type> The type of repository {hosted,external}
	[source] For external repositories, the source (e.g. registry.hub.docker.com/redis)

	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Invalid usage. sail repository add [<applicationName>/]<repositoryId> <type> [source]. Please see sail repository add --help")
		} else {
			source := ""
			if len(args) == 3 {
				source = args[2]
			}
			n := repositoryAddStruct{Type: args[1], Source: source}
			repositoryAdd(args[0], n)
		}
	},
}

type repositoryAddStruct struct {
	Type   string `json:"type"`
	Source string `json:"source,omitempty"`
}

func repositoryAdd(repositoryID string, args repositoryAddStruct) {
	// Split namespace and repository
	host, app, repo, tag, err := internal.ParseResourceName(repositoryID)
	internal.Check(err)

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	} else if len(tag) > 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid repository name. Please see sail repository add --help\n")
		os.Exit(1)
	}

	body, err := json.Marshal(args)
	internal.Check(err)

	path := fmt.Sprintf("/repositories/%s/%s", app, repo)
	internal.FormatOutputDef(internal.PostBodyWantJSON(path, body))

}
