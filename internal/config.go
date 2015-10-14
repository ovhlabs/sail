package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/docker/docker/cliconfig"
	"github.com/spf13/cobra"
)

type headers map[string]string

var (
	// Host points to the sailabove infrastructure wanted
	Host = os.Getenv("SAIL_HOST")
	// User of sailabove to use
	User = os.Getenv("SAIL_USER")
	// Password of sailabove account to use
	Password = os.Getenv("SAIL_PASSWORD")
	// ConfigDir points to the Docker configuration directory
	ConfigDir string
	// Verbose conditions the quantity of output of api requests
	Verbose bool
	// Format to use for output. One of 'json', 'yaml', 'pretty'
	Format string
	// Home fetches the user home directory
	Home = os.Getenv("HOME")
	// Headers to append to each requests. For debugging/internal purpose.
	Headers = make(headers)
)

func init() {
	if Host == "" {
		Host = "sailabove.io"
	}
	Cmd.AddCommand(cmdConfigShow)
}

// map --header arguments to map
func (h headers) Set(value string) error {
	chunks := strings.Split(value, "=")
	if len(chunks) != 2 {
		return fmt.Errorf("Invalid env var %s", value)
	}

	key := strings.TrimSpace(chunks[0])
	val := strings.TrimSpace(chunks[1])

	h[key] = val

	return nil
}

func (h headers) String() string {
	var i int
	var buffer = make([]string, len(h))

	for key, val := range h {
		buffer[i] = fmt.Sprintf("%s=%s", key, val)
		i++
	}
	return strings.Join(buffer, ", ")
}

func (h headers) Type() string {
	return fmt.Sprint("Headers")
}

// Cmd config
var Cmd = &cobra.Command{
	Use:     "config",
	Short:   "Config commands: sail config --help",
	Long:    `Config commands: sail config <command>`,
	Aliases: []string{"c"},
}

var cmdConfigShow = &cobra.Command{
	Use:   "show",
	Short: "Show Configuration: sail config show",
	Run: func(cmd *cobra.Command, args []string) {
		configShow()
	},
}

type configStruct struct {
	Username string `json:"username"`
	Host     string `json:"host"`
	Headers  string `json:"header,omitempty"`
}

func configShow() {
	var config configStruct

	ReadConfig()
	config.Username = User
	config.Host = Host
	config.Headers = Headers.String()

	data, err := json.Marshal(config)
	Check(err)

	FormatOutputDef(data)
}

// ReadConfig fetches docker config from ConfigDir
func ReadConfig() error {
	if !strings.Contains(Host, "://") {
		Host = "https://" + Host
	}

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

	url, err := url.ParseRequestURI(Host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid URL %s\n", Host)
		return err
	}

	for authHost, a := range c.AuthConfigs {

		if authHost == url.Host {
			if Verbose {
				fmt.Fprintf(os.Stderr, "Found in config file: Host %s Username:%s Password:<notShow>\n", authHost, a.Username)
			}

			if User == "" {
				User = a.Username
			}
			if Password == "" {
				Password = a.Password
			}

			if Verbose {
				fmt.Fprintf(os.Stderr, "Computed configuration: Host %s Username:%s Password:<notShow>\n", authHost, a.Username)
			}
			break
		}
	}

	if User == "" || Password == "" || Host == "" {
		return fmt.Errorf("Missing user, password or host in configuration. Did you forget to 'docker login %s' ?", url.Host)
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
