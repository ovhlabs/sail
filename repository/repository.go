package repository

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdRepositoryAdd)
	Cmd.AddCommand(cmdRepositoryDelete)
	Cmd.AddCommand(cmdRepositoryList)
}

// Cmd repository
var Cmd = &cobra.Command{
	Use:     "repository",
	Short:   "Repository commands: sail repository --help",
	Long:    `Repository commands: sail repository <command>`,
	Aliases: []string{"repo", "repositories"},
}
