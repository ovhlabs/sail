package compose

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/runabove/sail/internal"

	"github.com/spf13/cobra"
)

var (
	upFile    string
	upProject string
)

func cmdComposeUp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "sail compose up [<application>]",
		Run:   cmdUp,
	}

	wd, err := os.Getwd()
	if err == nil {
		t := strings.Split(wd, "/")
		wd = t[len(t)-1]
	}

	cmd.Flags().StringVarP(&upFile, "file", "", "docker-compose.yml", "Specify an alternate compose file")
	cmd.Flags().StringVarP(&upProject, "project-name", "p", wd, "Specify an alternate project name (default: directory name)")

	return cmd
}

func cmdUp(cmd *cobra.Command, args []string) {
	// FIXME: duplicate
	internal.ReadConfig()
	var ns string

	// Check args
	if len(args) > 1 {
		internal.Exit("Invalid usage. sail compose up [<application>]. Please see sail compose up -h\n")
	} else if len(args) > 1 {
		ns = args[0]
	} else {
		ns = internal.User
	}

	// Try to read file
	payload, err := ioutil.ReadFile(upFile)
	if err != nil {
		internal.Exit("Error reading compose file: %s\n", err)
	}

	// Execute request
	path := fmt.Sprintf("/applications/%s/fig/up?stream", ns)
	buffer, _, err := internal.Stream("POST", path, payload, internal.SetHeader("Content-Type", "application/x-yaml"))
	internal.Check(err)

	// Display api stream
	line, err := internal.DisplayStream(buffer)
	internal.Check(err)
	if line != nil {
		var data map[string]interface{}
		err = json.Unmarshal(line, &data)
		if err != nil {
			fmt.Printf("Error detected in API Return. Line: %s\n", line)
			return
		}

		fmt.Printf("Hostname: %v\n", data["hostname"])
		fmt.Printf("Running containers: %v/%v\n", data["container_number"], data["container_target"])
	}
}
