package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/cliconfig"
	"github.com/spf13/cobra"
)

var (
	// Host points to the sailabove infrastructure wanted
	Host string
	// User of sailabove to use
	User string
	// Password of sailabove account to use
	Password string
	// ConfigDir points to the Docker configuration directory
	ConfigDir string
	// Verbose conditions the quantity of output of api requests
	Verbose bool
	// Format to use for output. One of 'json', 'yaml', 'pretty'
	Format string
	// Home fetches the user home directory
	Home = os.Getenv("HOME")
)

func init() {
	Cmd.AddCommand(cmdConfigShow)
}

// Cmd config
var Cmd = &cobra.Command{
	Use:     "config",
	Short:   "Config commands : sail config --help",
	Long:    `Config commands : sail config <command>`,
	Aliases: []string{"c"},
}

var cmdConfigShow = &cobra.Command{
	Use:   "show",
	Short: "Show Configuration : sail config show",
	Run: func(cmd *cobra.Command, args []string) {
		configShow()
	},
}

type configStruct struct {
	Username string `json:"username"`
	Host     string `json:"host"`
}

func configShow() {
	var config configStruct

	ReadConfig()
	config.Username = User
	config.Host = Host

	data, err := json.Marshal(config)
	Check(err)

	FormatOutputDef(data)
}

// ReadConfig fetches docker config from ConfigDir
func ReadConfig() error {

	// if --user / --password are in args, take them.
	if User != "" && Password != "" {
		return nil
	}

	// otherwise try to take from docker config.json file
	c, err := cliconfig.Load(ConfigDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while reading config file in %s\n", ConfigDir)
		return err
	}

	if len(c.AuthConfigs) <= 0 {
		return fmt.Errorf("No Auth found in config file in %s", ConfigDir)
	}

	for authHost, a := range c.AuthConfigs {

		if authHost == Host {
			if Verbose {
				fmt.Fprintf(os.Stderr, "Found in config file : Host %s Username:%s Password:<notShow>\n", authHost, a.Username)
			}

			if User == "" {
				User = a.Username
			}
			if Password == "" {
				Password = a.Password
			}

			if Verbose {
				fmt.Fprintf(os.Stderr, "Computed configuration : Host %s Username:%s Password:<notShow>\n", authHost, a.Username)
			}
			break
		}
	}

	if User == "" || Password == "" || Host == "" {
		return fmt.Errorf("Invalid configuration, check user, password and host")
	}

	expandRegistryURL()
	return nil
}

func expandRegistryURL() {
	Host = Host + "/v1"
	if strings.HasPrefix(Host, "http") || strings.HasPrefix(Host, "https") {
		return
	}
	if ping("https://" + Host) {
		Host = "https://" + Host
		return
	}

	Host = "http://" + Host
	return
}

func ping(hostname string) bool {
	urlPing := hostname + "/_ping"
	if Verbose {
		fmt.Fprintf(os.Stderr, "Try ping on %s\n", urlPing)
	}
	req, _ := http.NewRequest("GET", urlPing, nil)
	initRequest(req)
	resp, err := getHTTPClient().Do(req)
	Check(err)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if Verbose {
			fmt.Fprintf(os.Stderr, "Ping OK on %s\n", urlPing)
		}
		return true
	}
	if Verbose {
		fmt.Fprintf(os.Stderr, "Ping KO on %s\n", urlPing)
	}
	return false
}
