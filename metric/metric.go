package metric

import "github.com/spf13/cobra"

func init() {
	cmdMetricToken.AddCommand(createCmd())
	cmdMetricToken.AddCommand(revokeCmd())

	Cmd.AddCommand(cmdMetricToken)
}

// Cmd application
var Cmd = &cobra.Command{
	Use:     "metric",
	Short:   "Metric commands: sail metric --help",
	Long:    `Metric commands: sail metric <command>`,
	Aliases: []string{"m", "metrics", "iot"},
}
