package domain

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(cmdDomainAttach)
	Cmd.AddCommand(cmdDomainDetach)
	Cmd.AddCommand(cmdDomainList)
}

// Cmd domain
var Cmd = &cobra.Command{
	Use:     "domain",
	Short:   "Service Domain commands : sail service domain --help",
	Long:    `Service Domain commands : sail service domain <command>`,
	Aliases: []string{"domains"},
}

type domainStruct struct {
	Pattern string `json:"pattern"`
	Method  string `json:"method"`
}
