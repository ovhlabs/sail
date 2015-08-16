package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"
)

func init() {
	cmdContainer.AddCommand(cmdContainerList)
	cmdContainer.AddCommand(cmdContainerInspect)

}

var cmdContainer = &cobra.Command{
	Use:     "container",
	Short:   "Container commands : sailgo container --help",
	Long:    `Container commands : sailgo container <command>`,
	Aliases: []string{"c", "containers"},
}

var cmdContainerList = &cobra.Command{
	Use:     "list",
	Short:   "List docker containers : sailgo container list [applicationName]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		containerList(getListApplications(args))
	},
}

var cmdContainerInspect = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a docker container : sailgo container inspect <applicationName> <containerId>",
	Long: `Inspect a docker container : sailgo container inspect <applicationName> <containerId>
	\"example : sailgo container inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Invalid usage. sailgo container inspect <applicationName> <containerId>. Please see sailgo container inspect --help")
		} else {
			fmt.Println(getWantJSON(fmt.Sprintf("/applications/%s/containers/%s", args[0], args[1])))
		}
	},
}

func containerList(apps []string) {

	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	titles := []string{"APPLICATION", "SERVICE", "CONTAINER", "STATE", "DEPLOYED"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	containers := []string{}
	var container map[string]interface{}
	for _, app := range apps {
		b := reqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/containers", app), nil)
		check(json.Unmarshal(b, &containers))
		for _, containerID := range containers {
			b := reqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/containers/%s", app, containerID), nil)
			check(json.Unmarshal(b, &container))
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", app, container["service"], container["name"], strings.ToUpper(container["state"].(string)), container["deployment_date"])
			w.Flush()
		}
	}
}
