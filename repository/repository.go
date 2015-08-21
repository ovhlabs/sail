package repository

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdRepositoryList)

	// TODO
	//sail repository add    Add a new docker repository
	//sail repository rm     Delete a repository
}

// Cmd repository
var Cmd = &cobra.Command{
	Use:     "repository",
	Short:   "Repository commands : sailgo repository --help",
	Long:    `Repository commands : sailgo repository <command>`,
	Aliases: []string{"repo", "repositories"},
}
