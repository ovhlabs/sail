package network

import (
	"encoding/json"
	"fmt"
	"strings"

	"stash.ovh.net/sailabove/sail/internal"

	"github.com/spf13/cobra"
)

var cmdNetworkAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a new private network : sail network add <applicationName>/<networkId> subnet",
	Long:  `Add a new private network : sail network add <applicationName>/<networkId> subnet`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Invalid usage. sail network add <applicationName>/<networkId> subnet. Please see sail network add --help")
		} else {
			n := networkAddStruct{Subnet: args[1]}
			networkAdd(args[0], n)
		}
	},
}

type networkAddStruct struct {
	Subnet string `json:"subnet"`
}

func networkAdd(networkID string, args networkAddStruct) {
	t := strings.Split(networkID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sail network add <applicationName>/<networkId>. Please see sail network add --help")
		return
	}

	body, err := json.Marshal(args)
	internal.Check(err)

	path := fmt.Sprintf("/applications/%s/networks/%s", t[0], t[1])
	fmt.Println(internal.PostBodyWantJSON(path, body))

}
