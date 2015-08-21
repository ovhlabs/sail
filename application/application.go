package application

import (
	"fmt"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sailgo/internal"
)

func init() {
	Cmd.AddCommand(cmdApplicationList)
	Cmd.AddCommand(cmdApplicationInspect)

	cmdApplicationDomain.AddCommand(cmdApplicationDomainList)
	cmdApplicationDomain.AddCommand(cmdApplicationDomainAttach)
	cmdApplicationDomain.AddCommand(cmdApplicationDomainDetach)
	Cmd.AddCommand(cmdApplicationDomain)
}

var Cmd = &cobra.Command{
	Use:     "application",
	Short:   "Application commands : sailgo application --help",
	Long:    `Application commands : sailgo application <command>`,
	Aliases: []string{"a", "app", "apps", "applications"},
}

var cmdApplicationList = &cobra.Command{
	Use:     "list",
	Short:   "List granted apps : sailgo application list",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(internal.GetWantJSON("/applications"))
	},
}

var cmdApplicationInspect = &cobra.Command{
	Use:   "inspect",
	Short: "Details of an app : sailgo application inspect <applicationName>",
	Long: `Details of an app : sailgo application inspect <applicationName>
	\"example : sailgo application inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Println("Invalid usage. Please see sailgo application inspect --help")
		} else {
			fmt.Println(internal.GetWantJSON(fmt.Sprintf("/applications/%s", args[0])))
		}
	},
}

var cmdApplicationDomain = &cobra.Command{
	Use:     "domain",
	Short:   "Application Domain commands : sailgo application domain --help",
	Long:    `Application Domain commands : sailgo application domain <command>`,
	Aliases: []string{"domains"},
}

var cmdApplicationDomainList = &cobra.Command{
	Use:     "list",
	Short:   "List domains and routes on the HTTP load balancer : sailgo application domain list <applicationName>",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Println("Invalid usage. Please see sailgo application domain list --help")
		} else {
			// cmdApplicationDomainList TODO ? Tab view with headers ['DOMAIN', 'SERVICE', 'METHOD', 'PATTERN']
			fmt.Println(internal.GetWantJSON(fmt.Sprintf("/applications/%s/attached-domains", args[0])))
		}
	},
}

var cmdApplicationDomainAttach = &cobra.Command{
	Use:     "attach",
	Short:   "Attach a domain on the HTTP load balancer : sailgo application domain attach <applicationName> <domainName>",
	Aliases: []string{"add"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Invalid usage. Please see sailgo application domain attach --help")
		} else {
			fmt.Println(internal.PostWantJSON(fmt.Sprintf("/applications/%s/attached-domains/%s", args[0], args[1])))
		}
	},
}

var cmdApplicationDomainDetach = &cobra.Command{
	Use:     "detach",
	Short:   "Detach a domain from the HTTP load balancer : sailgo application domain detach <applicationName> <domainName>",
	Aliases: []string{"add"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Invalid usage. Please see sailgo application domain attach --help")
		} else {
			fmt.Println(internal.DeleteWantJSON(fmt.Sprintf("/applications/%s/attached-domains/%s", args[0], args[1])))
		}
	},
}
