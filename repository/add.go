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
	Short: "Add an external repository: sail repository add [<applicationName>/]<repositoryName> [<registryURL>] [<namespace>/]<externalRepositoryName> ",
	Long: `Add an external repository: sail repository add [<applicationName>/]<repositoryName> [<registryURL>] [<user>/]<externalRepositoryName>[:<tag>]
examples:
	- sail repository add myapp/myrepo docker.private.registry.com privateRegistryAccount/mysql:5.5
	- sail repository add myapp/myrepo dockerHubAccount/mysql:5.5
	- sail repository add myapp/myrepo mysql

[<user>] The account which hold the distant repository. If ommitted it use official docker image.
<externalRepositoryName> the name of the distant repository. Example : tutum/nginx
[<registryURL>] Url of the registry. Default value : the docker hub (e.g. https://hub.docker.com)
Only pubic repositories are supported yet. ADD NOTE FOR PRIVATE REGISTRIES
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var registryURL, externalRepositoryName, repositoryName string
		switch len(args) {
		case 2:
			repositoryName = args[0]
			registryURL = ""
			externalRepositoryName = args[1]
		case 3:
			repositoryName = args[0]
			registryURL = args[1]
			externalRepositoryName = args[2]

		default:
			fmt.Fprintln(os.Stderr, "Invalid usage. sail repository add [<applicationName>/]<repositoryName> [<registryURL>/][<user>/]<externalRepositoryName>[:<tag>]. Please see sail repository add --help")
			return
		}

		n := repositoryAddStruct{
			RegistryURL:            registryURL,
			ExternalRepositoryName: externalRepositoryName,
		}
		repositoryAdd(repositoryName, n)
	},
}

type repositoryAddStruct struct {
	RegistryURL            string `json:"registryURL,omitempty"`
	ExternalRepositoryName string `json:"externalRepositoryName"`
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
