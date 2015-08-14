package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/viper"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"
)

var verbose bool
var host, user, password, configDir string
var home = os.Getenv("HOME")

var rootCmd = &cobra.Command{
	Use:   "sailgo",
	Short: "Sailabove - Command Line Tool",
	Long:  `Sailabove - Command Line Tool`,
}

func main() {
	addCommands()
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "", "Docker index host, facultative if you have a "+home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "Docker index user, facultative if you have a "+home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Docker index password, facultative if you have a "+home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&configDir, "configDir", "", home+"/.docker", "configuration directory, default is "+home+"/.docker/")

	rootCmd.Execute()
}

// AddCommands adds child commands to the root command rootCmd.
func addCommands() {
	rootCmd.AddCommand(cmdConfig)
	rootCmd.AddCommand(cmdCompose)
	rootCmd.AddCommand(cmdContainer)
	rootCmd.AddCommand(cmdMe)
	rootCmd.AddCommand(cmdNetwork)
	rootCmd.AddCommand(cmdRepository)
	rootCmd.AddCommand(cmdService)
	rootCmd.AddCommand(cmdVersion)
}

func initRequest(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func getHTTPClient() *http.Client {
	tr := &http.Transport{}
	return &http.Client{Transport: tr}
}

func reqWant(method string, wantCode int, path string, jsonStr []byte) []byte {

	readConfig()

	// TODO URL
	/*if viper.GetString("host") == "" {
		fmt.Println("Invalid Configuration : invalid URL. See sailgo config --help")
		os.Exit(1)
	}*/

	var req *http.Request
	if jsonStr != nil {
		req, _ = http.NewRequest(method, viper.GetString("host")+path, bytes.NewReader(jsonStr))
	} else {
		req, _ = http.NewRequest(method, viper.GetString("host")+path, nil)
	}

	initRequest(req)
	resp, err := getHTTPClient().Do(req)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode != wantCode || verbose {
		fmt.Printf("Response Status:%s\n", resp.Status)
		fmt.Printf("Request path :%s\n", viper.GetString("host")+path)
		fmt.Printf("Request :%s\n", string(jsonStr))
		fmt.Printf("Response Headers:%s\n", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Response Body:%s\n", string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	return body
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
