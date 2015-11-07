package update

import "github.com/spf13/cobra"

// used by CI to inject url for downloading with sail update.
// value of urlUpdate injected at build time
// full URL update is constructed with architecture var :
// urlUpdate + "sail-" + architecture, sail is the binary
var urlUpdateSnapshot string

func init() {
	if urlUpdateSnapshot != "" {
		Cmd.AddCommand(cmdUpdateSnapshot)
	}
}

var cmdUpdateSnapshot = &cobra.Command{
	Use:     "snapshot",
	Short:   "Update sail to latest snapshot version: sail update snapshot",
	Long:    `sail update snapshot`,
	Aliases: []string{"snap"},
	Run: func(cmd *cobra.Command, args []string) {
		doUpdate(urlUpdateSnapshot, architecture)
	},
}
