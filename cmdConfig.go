package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/viper"
)

func init() {
	cmdConfig.AddCommand(cmdConfigTemplate)
	cmdConfig.AddCommand(cmdConfigShow)
}

var cmdConfig = &cobra.Command{
	Use:     "config",
	Short:   "Config commands : sailgo config --help",
	Long:    `Config commands : sailgo config <command>`,
	Aliases: []string{"c"},
}

var cmdConfigTemplate = &cobra.Command{
	Use:   "template",
	Short: "Write a template configuration file in $HOME/.sailgo/config.json : sailgo config template",
	Run: func(cmd *cobra.Command, args []string) {
		writeTemplate()
	},
}

var cmdConfigShow = &cobra.Command{
	Use:   "show",
	Short: "Show Configuration : sailgo config show",
	Run: func(cmd *cobra.Command, args []string) {
		show()
	},
}

type templateJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
}

func writeTemplate() {
	var templateJSON templateJSON

	if viper.GetString("username") != "" {
		templateJSON.Username = viper.GetString("username")
	}
	if viper.GetString("password") != "" {
		templateJSON.Password = viper.GetString("password")
	}
	if viper.GetString("url") != "" {
		templateJSON.URL = viper.GetString("url")
	}
	jsonStr, err := json.MarshalIndent(templateJSON, "", "  ")
	check(err)
	jsonStr = append(jsonStr, '\n')
	filename := os.Getenv("HOME") + "/.sailgo/config.json"
	check(ioutil.WriteFile(filename, jsonStr, 0600))
	fmt.Printf("%s is written\n", filename)
}

func show() {
	readConfig()
	fmt.Printf("username:%s\n", viper.GetString("username"))
	fmt.Printf("password:%s\n", viper.GetString("password"))
	fmt.Printf("url:%s\n", viper.GetString("url"))
}
