package metric

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdMetricToken = &cobra.Command{
	Use:     "token",
	Short:   "Metric token commands: sail metric token --help",
	Long:    `Metric token commands: sail metric token <command>`,
	Aliases: []string{"metrics"},
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Metrics token for a given application: sail metric token create <applicationName>",
		Aliases: []string{"create", "new", "c"},
		Run:     cmdCreate,
	}

	return cmd
}

func revokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "revoke",
		Short:   "Revoke metrics token for a given application: sail application metric revoke <applicationName> <token-username>",
		Aliases: []string{"delete", "del", "r"},
		Run:     cmdRevoke,
	}

	return cmd
}

func cmdCreate(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail metric token create <applicationName>. Please see sail metric token create --help\n"
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	// Get args
	application := args[0]
	path := fmt.Sprintf("/applications/%s/metrics/token", application)
	internal.FormatOutputDef(internal.PostWantJSON(path))
}

func cmdRevoke(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail metric token revoke <applicationName> <token>. Please see sail metric token revoke --help\n"
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	// Get args
	application := args[0]
	token := args[1]

	path := fmt.Sprintf("/applications/%s/metrics/token/%s", application, token)
	data := internal.DeleteWantJSON(path)

	internal.FormatOutput(data, func(data []byte) {
		fmt.Fprintf(os.Stderr, "Disabled IoT token %s from service %s\n", token, application)
	})
}
