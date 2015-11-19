package application

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdApplicationList)
	Cmd.AddCommand(cmdApplicationShow)

	cmdApplicationDomain.AddCommand(cmdApplicationDomainList)
	cmdApplicationDomain.AddCommand(cmdApplicationDomainDetach)

	Cmd.AddCommand(cmdApplicationDomain)

	cmdApplicationWebhook.AddCommand(cmdApplicationWebhookList)
	cmdApplicationWebhook.AddCommand(cmdApplicationWebhookAdd)
	cmdApplicationWebhook.AddCommand(cmdApplicationWebhookDelete)

	Cmd.AddCommand(cmdApplicationWebhook)

	cmdApplicationEnv.AddCommand(cmdApplicationListEnv)
	cmdApplicationEnv.AddCommand(cmdApplicationSetEnv)
	cmdApplicationEnv.AddCommand(cmdApplicationDelEnv)

	Cmd.AddCommand(cmdApplicationEnv)
}

// Cmd application
var Cmd = &cobra.Command{
	Use:     "application",
	Short:   "Application commands: sail application --help",
	Long:    `Application commands: sail application <command>`,
	Aliases: []string{"a", "app", "apps", "applications"},
}
