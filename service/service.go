package service

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdServiceAttach)
	Cmd.AddCommand(cmdServiceList)
	Cmd.AddCommand(cmdServiceInspect)
	Cmd.AddCommand(redeployCmd())
	Cmd.AddCommand(addCmd())
	Cmd.AddCommand(rmCmd())
	//TODO
	// sail services logs           Fetch the logs of a service
	// sail services stop           Stop a docker service
	// sail services start          Start a docker service
	// sail services scale          Scale a docker service
	// sail services domain-list    List domains on the HTTP load balancer
	// sail services domain-attach  Attach a domain on the HTTP load balancer
	// sail services domain-detach  Detach a domain from the HTTP load balancer

}

// Cmd service
var Cmd = &cobra.Command{
	Use:     "service",
	Short:   "Service commands : sailgo service --help",
	Long:    `Service commands : sailgo service <command>`,
	Aliases: []string{"services"},
}
