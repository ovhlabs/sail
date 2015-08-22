package repository

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdRepositoryAdd)
	Cmd.AddCommand(cmdRepositoryRm)
	Cmd.AddCommand(cmdRepositoryList)
}

// Cmd repository
var Cmd = &cobra.Command{
	Use:     "repository",
	Short:   "Repository commands : sailgo repository --help",
	Long:    `Repository commands : sailgo repository <command>`,
	Aliases: []string{"repo", "repositories"},
}
