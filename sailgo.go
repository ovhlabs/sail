package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"
	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/viper"
)

var verbose bool
var url, auth, email, configFile string
var home = os.Getenv("HOME")

var rootCmd = &cobra.Command{
	Use:   "sailgo",
	Short: "Sailabove - Command Line Tool",
	Long:  `Sailabove - Command Line Tool`,
}

func main() {
	addCommands()
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&auth, "auth", "u", "", "auth, facultative if you have a "+home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&email, "email", "p", "", "email, facultative if you have a "+home+"/.docker/config.json file")
	rootCmd.PersistentFlags().StringVarP(&configFile, "configFile", "c", home+"/.docker/config.json", "configuration file, default is "+home+"/.docker/config.json")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("auth", rootCmd.PersistentFlags().Lookup("auth"))
	viper.BindPFlag("email", rootCmd.PersistentFlags().Lookup("email"))

	rootCmd.Execute()
}

func readConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		viper.ReadInConfig() // Find and read the config file
		if verbose {
			fmt.Printf("Using config file %s\n", configFile)
		}
	}
}

//AddCommands adds child commands to the root command rootCmd.
func addCommands() {
	rootCmd.AddCommand(cmdMe)
	rootCmd.AddCommand(cmdConfig)
	rootCmd.AddCommand(cmdVersion)
}

func getSkipLimit(tail []string) (string, string) {
	skip := "0"
	limit := "10"
	if len(tail) == 3 {
		skip = tail[1]
		limit = tail[2]
	} else if len(tail) == 2 {
		skip = tail[0]
		limit = tail[1]
	}
	return skip, limit
}

func initRequest(req *http.Request) {
	// TODOs
	// req.Header.Set("Content-Type", "application/json")
}

func getHTTPClient() *http.Client {
	tr := &http.Transport{}
	return &http.Client{Transport: tr}
}

func reqWant(method string, wantCode int, path string, jsonStr []byte) []byte {

	readConfig()

	// TODO URL
	/*if viper.GetString("url") == "" {
		fmt.Println("Invalid Configuration : invalid URL. See sailgo config --help")
		os.Exit(1)
	}*/

	var req *http.Request
	if jsonStr != nil {
		req, _ = http.NewRequest(method, viper.GetString("url")+path, bytes.NewReader(jsonStr))
	} else {
		req, _ = http.NewRequest(method, viper.GetString("url")+path, nil)
	}

	initRequest(req)
	resp, err := getHTTPClient().Do(req)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode != wantCode || verbose {
		fmt.Printf("Response Status:%s\n", resp.Status)
		fmt.Printf("Request path :%s\n", viper.GetString("url")+path)
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
