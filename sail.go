package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/runabove/sail/application"
	"github.com/runabove/sail/compose"
	"github.com/runabove/sail/container"
	"github.com/runabove/sail/internal"
	"github.com/runabove/sail/me"
	"github.com/runabove/sail/metric"
	"github.com/runabove/sail/network"
	"github.com/runabove/sail/repository"
	"github.com/runabove/sail/service"
	"github.com/runabove/sail/update"
	"github.com/runabove/sail/version"
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
	rootCmd.PersistentFlags().StringVarP(&internal.Host, "host", "H", internal.Host, "Docker index host [$SAIL_HOST], optional if you have a "+internal.Home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&internal.User, "username", "U", internal.User, "Docker index user [$SAIL_USER], optional if you have a "+internal.Home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&internal.Password, "password", "P", internal.Password, "Docker index password [$SAIL_PASSWORD], optional if you have a "+internal.Home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&internal.ConfigDir, "configDir", "", internal.Home+"/.docker", "configuration directory, default is "+internal.Home+"/.docker/")
	rootCmd.PersistentFlags().Var(&internal.Headers, "header", "'KEY=value' headers to append to each requests. For debugging/internal purpose.")
	rootCmd.Execute()
}

// AddCommands adds child commands to the root command rootCmd.
func addCommands() {
	rootCmd.AddCommand(application.Cmd)
	rootCmd.AddCommand(compose.Cmd)
	rootCmd.AddCommand(internal.Cmd)
	rootCmd.AddCommand(container.Cmd)
	rootCmd.AddCommand(me.Cmd)
	rootCmd.AddCommand(metric.Cmd)
	rootCmd.AddCommand(network.Cmd)
	rootCmd.AddCommand(repository.Cmd)
	rootCmd.AddCommand(service.Cmd)
	rootCmd.AddCommand(update.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(autocompleteCmd)
}

var autocompleteCmd = &cobra.Command{
	Use:   "autocomplete <path>",
	Short: "Generate bash autocompletion file for sail",
	Long:  `Generate bash autocompletion file for sail`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "Wrong usage: sail autocomplete <path>\n")
			return
		}
		rootCmd.GenBashCompletionFile(args[0])
		fmt.Fprintf(os.Stderr, "Completion file generated.\n")
		fmt.Fprintf(os.Stderr, "You may now run `source %s`\n", args[0])
	},
}
