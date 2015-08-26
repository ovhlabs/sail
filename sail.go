package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sail/application"
	"stash.ovh.net/sailabove/sail/compose"
	"stash.ovh.net/sailabove/sail/container"
	"stash.ovh.net/sailabove/sail/internal"
	"stash.ovh.net/sailabove/sail/me"
	"stash.ovh.net/sailabove/sail/network"
	"stash.ovh.net/sailabove/sail/repository"
	"stash.ovh.net/sailabove/sail/service"
	"stash.ovh.net/sailabove/sail/version"
)

var rootCmd = &cobra.Command{
	Use:   "sail",
	Short: "Sailabove - Command Line Tool",
	Long:  `Sailabove - Command Line Tool`,
}

func main() {
	addCommands()
	rootCmd.PersistentFlags().BoolVarP(&internal.Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&internal.Format, "format", "f", "pretty", "choose format output. One of 'json', 'yaml' and 'pretty'")
	rootCmd.PersistentFlags().StringVarP(&internal.Host, "host", "H", "sailabove.io", "Docker index host, optional if you have a "+internal.Home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&internal.User, "user", "U", "", "Docker index user, optional if you have a "+internal.Home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&internal.Password, "password", "P", "", "Docker index password, optional if you have a "+internal.Home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&internal.ConfigDir, "configDir", "", internal.Home+"/.docker", "configuration directory, default is "+internal.Home+"/.docker/")
	rootCmd.Execute()
}

// AddCommands adds child commands to the root command rootCmd.
func addCommands() {
	rootCmd.AddCommand(application.Cmd)
	rootCmd.AddCommand(compose.Cmd)
	rootCmd.AddCommand(internal.Cmd)
	rootCmd.AddCommand(container.Cmd)
	rootCmd.AddCommand(me.Cmd)
	rootCmd.AddCommand(network.Cmd)
	rootCmd.AddCommand(repository.Cmd)
	rootCmd.AddCommand(service.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(autocompleteCmd)
}

var autocompleteCmd = &cobra.Command{
	Use:   "autocomplete <path>",
	Short: "Generate bash autocompletion file for sail",
	Long:  `Generate bash autocompletion file for sail`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Wrong usage: sail autocomplete <path>\n")
			return
		}
		rootCmd.GenBashCompletionFile(args[0])
		fmt.Printf("Completion file generated.\n")
		fmt.Printf("You may now run `source %s`\n", args[0])
	},
}
