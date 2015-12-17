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
	Short: "Add an external repository: sail repository add [<applicationName>/]<repositoryName> [<source>/][<namespace>/]<externalRepositoryName> ",
	Long: `Add an external repository: sail repository add [<applicationName>/]<repositoryName> [<source>/][<user>/]<externalRepositoryName>
examples:
	- sail repository add myapp/myrepo docker.private.registry/privateRegistryAccount/mysql
	- sail repository add myapp/myrepo dockerHubAccount/mysql
	- sail repository add myapp/myrepo mysql

[<user>] The account which hold the distant repository. If ommitted it use official docker image.
<externalRepositoryName> the name of the distant repository. Example : tutum/nginx
[<source>] For external repositories. Default value : the docker hub (e.g. https://hub.docker.com)
Only pubic repositories are supported yet. ADD NOTE FOR PRIVATE REGISTRIES
	`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 2:
			repositoryName := args[0]
			source := args[1]

		default:
			fmt.Fprintln(os.Stderr, "Invalid usage. sail repository add [<applicationName>/]<repositoryName> [<source>/][<user>/]<externalRepositoryName>. Please see sail repository add --help")
			return
		}

		n := repositoryAddStruct{Source: source}
		repositoryAdd(repositoryName, n)
	},
}

type repositoryAddStruct struct {
	Source string `json:"source"`
}

func repositoryAdd(repositoryName string, args repositoryAddStruct) {
	// Split namespace and repository
	host, app, repo, tag, err := internal.ParseResourceName(repositoryName)
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
