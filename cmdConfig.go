package main

import (
	"fmt"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/docker/docker/cliconfig"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/viper"
)

func init() {
	cmdConfig.AddCommand(cmdConfigShow)
}

var cmdConfig = &cobra.Command{
	Use:     "config",
	Short:   "Config commands : sailgo config --help",
	Long:    `Config commands : sailgo config <command>`,
	Aliases: []string{"c"},
}

var cmdConfigShow = &cobra.Command{
	Use:   "show",
	Short: "Show Configuration : sailgo config show",
	Run: func(cmd *cobra.Command, args []string) {
		show()
	},
}

func show() {
	readConfig()
	fmt.Printf("username:%s\n", viper.GetString("username"))
	fmt.Printf("password:%s\n", viper.GetString("password"))
	fmt.Printf("url:%s\n", viper.GetString("url"))
}

func readConfig() error {
	//TODO put -H ServerAddress credentials in a struct
	c, err := cliconfig.Load(configDir)
	fmt.Printf("Config file : %+v", c)
	return err
}
