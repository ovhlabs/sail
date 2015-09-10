package application

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sail/internal"
)

var (
	tokenBody Token
)

var cmdApplicationMetric = &cobra.Command{
	Use:     "metric",
	Short:   "Application Metric commands: sail application metric --help",
	Long:    `Application Metric commands: sail application metric <command>`,
	Aliases: []string{"metrics"},
}

func tokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token",
		Short:   "Metrics token for a given application: sail application metric token <applicationName>",
		Aliases: []string{"getToken"},
		Run:     cmdToken,
	}

	cmd.Flags().IntVarP(&tokenBody.Validity, "validity", "", 86400, "Validity in seconds. Defaults to a full day")

	return cmd
}

func revokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "revoke",
		Short:   "Revoke metrics token for a given application: sail application metric revoke <applicationName> <token>",
		Aliases: []string{"revokeToken"},
		Run:     cmdRevoke,
	}

	return cmd
}

// Token struct holds all parameters sent to GET /applications/%s/metrics/token
type Token struct {
	Validity int `json:"validity"`
}

func cmdToken(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail application metric token <applicationName>. Please see sail application metric token --help\n"
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	// Get args
	application := args[0]
	body, err := json.Marshal(tokenBody)
	internal.Check(err)

	path := fmt.Sprintf("/applications/%s/metrics/token", application)
	internal.FormatOutputDef(internal.PostBodyWantJSON(path, body))
}

func cmdRevoke(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail application metric revoke <applicationName> <token>. Please see sail application metric revoke --help\n"
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	// Get args
	application := args[0]
	token := args[1]

	path := fmt.Sprintf("/applications/%s/metrics/token/%s", application, token)
	internal.FormatOutputDef(internal.DeleteWantJSON(path))
}
