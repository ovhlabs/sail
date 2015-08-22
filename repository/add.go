package repository

import (
	"encoding/json"
	"fmt"
	"strings"

	"stash.ovh.net/sailabove/sailgo/internal"

	"github.com/spf13/cobra"
)

var cmdRepositoryAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a new repository : sailgo repository add <applicationName>/<repositoryId> <type> [source]",
	Long: `Add a new repository : sailgo repository add <applicationName>/<repositoryId> <type> [source]

	<type> The type of repository {hosted,external}
	[source] For external repositories, the source (e.g. registry.hub.docker.com/redis)

	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Invalid usage. sailgo repository add <applicationName>/<repositoryId> <type> [source]. Please see sailgo repository add --help")
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
	t := strings.Split(repositoryID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sailgo repository add <applicationName>/<repositoryId>. Please see sailgo repository add --help")
		return
	}

	body, err := json.Marshal(args)
	internal.Check(err)

	path := fmt.Sprintf("/repositories/%s/%s", t[0], t[1])
	fmt.Println(internal.PostBodyWantJSON(path, body))

}
