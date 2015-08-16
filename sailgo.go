package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"
)

var verbose, pretty bool
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
	rootCmd.PersistentFlags().BoolVarP(&pretty, "pretty", "t", false, "Pretty Print Json Output")
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "sailabove.io", "Docker index host, facultative if you have a "+home+"/.docker/config.json file")
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
	req.Header.Set("User-Agent", "Sailabove sailgo CLI/"+VERSION)
}

func getHTTPClient() *http.Client {
	tr := &http.Transport{}
	return &http.Client{Transport: tr}
}

func getWantReturn(path string) string {
	return getJSON(reqWant("GET", http.StatusOK, path, nil))
}

func reqWantJSON(method string, wantCode int, path string, body []byte) string {
	return getJSON(reqWant(method, wantCode, path, body))
}

func reqWant(method string, wantCode int, path string, jsonStr []byte) []byte {

	readConfig()

	var req *http.Request
	if jsonStr != nil {
		req, _ = http.NewRequest(method, host+path, bytes.NewReader(jsonStr))
	} else {
		req, _ = http.NewRequest(method, host+path, nil)
	}

	initRequest(req)
	req.SetBasicAuth(user, password)
	resp, err := getHTTPClient().Do(req)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode != wantCode || verbose {
		fmt.Printf("Response Status : %s\n", resp.Status)
		fmt.Printf("Request path : %s\n", host+path)
		fmt.Printf("Request Headers : %s\n", req.Header)
		fmt.Printf("Request Body : %s\n", string(jsonStr))
		fmt.Printf("Response Headers : %s\n", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Response Body : %s\n", string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	return body
}

func getJSON(s []byte) string {
	if pretty {
		var out bytes.Buffer
		json.Indent(&out, s, "", "  ")
		return out.String()
	}
	return string(s)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
