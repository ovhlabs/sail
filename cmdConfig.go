package main

import (
	"fmt"
	"net/http"
	"strings"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/docker/docker/cliconfig"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"
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
		configShow()
	},
}

func configShow() {
	readConfig()
	fmt.Printf("username:%s\n", user)
	fmt.Printf("host:%s\n", host)
}

func readConfig() error {

	// if --user / --password are in args, take them.
	if user != "" && password != "" {
		return nil
	}

	// otherwise try to take from docker config.json file
	c, err := cliconfig.Load(configDir)
	if err != nil {
		fmt.Printf("Error while reading config file in %s\n", configDir)
		return err
	}

	if len(c.AuthConfigs) <= 0 {
		return fmt.Errorf("No Auth found in config file in %s", configDir)
	}

	for authHost, a := range c.AuthConfigs {
		if authHost == host {
			if verbose {
				fmt.Printf("Found in config file : Host %s Username:%s Password:<notShow>\n", authHost, a.Username)
			}

			if user == "" {
				user = a.Username
			}
			if password == "" {
				password = a.Password
			}

			if verbose {
				fmt.Printf("Computed configuration : Host %s Username:%s Password:<notShow>\n", authHost, a.Username)
			}
			break
		}
	}

	if user == "" || password == "" || host == "" {
		return fmt.Errorf("Invalid configuration, check user, password and host")
	}

	expandRegistryURL()
	return nil
}

func expandRegistryURL() {
	host = host + "/v1"
	if strings.HasPrefix(host, "http") || strings.HasPrefix(host, "https") {
		return
	}
	if ping("https://" + host) {
		host = "https://" + host
		return
	}

	host = "http://" + host
	return
}

func ping(hostname string) bool {
	urlPing := hostname + "/_ping"
	if verbose {
		fmt.Printf("Try ping on %s\n", urlPing)
	}
	req, _ := http.NewRequest("GET", urlPing, nil)
	initRequest(req)
	resp, err := getHTTPClient().Do(req)
	check(err)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if verbose {
			fmt.Printf("Ping OK on %s\n", urlPing)
		}
		return true
	}
	if verbose {
		fmt.Printf("Ping KO on %s\n", urlPing)
	}
	return false
}
