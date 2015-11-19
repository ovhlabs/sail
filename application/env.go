package application

import (
	"encoding/json"
	"fmt"
	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type env struct {
	Data string `json:"data"`
}

var cmdApplicationEnv = &cobra.Command{
	Use:   "env",
	Short: "Application env commands: sail application env --help",
	Long:  `Application env commands: sail application env <command>`,
}

var cmdApplicationListEnv = &cobra.Command{
	Use:     "list",
	Short:   "List environment variables of given application: sail application env list [<applicationName>]",
	Aliases: []string{},
	Run:     cmdListEnv,
}

var cmdApplicationSetEnv = &cobra.Command{
	Use:     "set",
	Short:   "Set an environment variable for given application: sail application env set [<applicationName>] <KEY=VALUE>",
	Aliases: []string{"add"},
	Run:     cmdSetEnv,
}

var cmdApplicationDelEnv = &cobra.Command{
	Use:     "delete",
	Short:   "Delete an environment variable for given application: sail application env delete [<applicationName>] <KEY>",
	Aliases: []string{"del", "remove", "rm"},
	Run:     cmdDelEnv,
}

func cmdListEnv(cmd *cobra.Command, args []string) {
	var applicationName string

	switch len(args) {
	case 0:
		applicationName = internal.GetUserName()
	case 1:
		applicationName = args[0]
	default:
		fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application env list --help")
		return
	}

	internal.FormatOutputDef(internal.GetWantJSON(fmt.Sprintf("/applications/%s/env", applicationName)))
}

func cmdSetEnv(cmd *cobra.Command, args []string) {
	var cmdEnvBody env
	var applicationName string
	var parsedData []string

	switch len(args) {
	case 1:
		applicationName = internal.GetUserName()
		parsedData = strings.SplitN(args[0], "=", 2)
	case 2:
		applicationName = args[0]
		parsedData = strings.SplitN(args[1], "=", 2)
	default:
		fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application env set --help")
		return
	}

	cmdEnvBody.Data = parsedData[1]
	jsonStr, err := json.Marshal(cmdEnvBody)
	internal.Check(err)
	internal.FormatOutputDef(internal.PostBodyWantJSON(fmt.Sprintf("/applications/%s/env/%s", applicationName, parsedData[0]), jsonStr))
}

func cmdDelEnv(cmd *cobra.Command, args []string) {
	var applicationName string
	var key string

	switch len(args) {
	case 1:
		applicationName = internal.GetUserName()
		key = args[0]
	case 2:
		applicationName = args[0]
		key = args[1]
	default:
		fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application env set --help")
		return
	}

	internal.FormatOutputDef(internal.DeleteWantJSON(fmt.Sprintf("/applications/%s/env/%s", applicationName, key)))
}
