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
	cmdService.AddCommand(cmdServiceList)
	cmdService.AddCommand(cmdServiceInspect)

	// TODO
	// sail services add            Add a new docker service
	// sail services rm             Delete a docker service
	// sail services attach         Attach to the console of the service containers
	// sail services logs           Fetch the logs of a service
	// sail services redeploy       Redeploy a docker service
	// sail services stop           Stop a docker service
	// sail services start          Start a docker service
	// sail services scale          Scale a docker service
	// sail services domain-list    List domains on the HTTP load balancer
	// sail services domain-attach  Attach a domain on the HTTP load balancer
	// sail services domain-detach  Detach a domain from the HTTP load balancer

}

var cmdService = &cobra.Command{
	Use:     "service",
	Short:   "Service commands : sailgo service --help",
	Long:    `Service commands : sailgo service <command>`,
	Aliases: []string{"services"},
}

var cmdServiceList = &cobra.Command{
	Use:     "list",
	Short:   "List the docker services : sailgo service list [applicationName]",
	Aliases: []string{"ls", "ps"},
	Run: func(cmd *cobra.Command, args []string) {
		serviceList(getListApplications(args))
	},
}

var cmdServiceInspect = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect a docker service : sailgo service inspect <applicationName>/<serviceId>",
	Long: `Inspect a docker service : sailgo service inspect <applicationName>/<serviceId>
	\"example : sailgo service inspect myApp"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Invalid usage. sailgo service inspect <applicationName>/<serviceId>. Please see sailgo service inspect --help")
		} else {
			serviceInspect(args[0])
		}
	},
}

func serviceInspect(serviceID string) {
	t := strings.Split(serviceID, "/")
	if len(t) != 2 {
		fmt.Println("Invalid usage. sailgo service inspect <applicationName>/<serviceId>. Please see sailgo service inspect --help")
	} else {
		fmt.Println(getWantJSON(fmt.Sprintf("/applications/%s/services/%s", t[0], t[1])))
	}
}

func serviceList(apps []string) {
	w := tabwriter.NewWriter(os.Stdout, 27, 1, 2, ' ', 0)
	titles := []string{"NAME", "REPOSITORY", "IMAGE ID", "STATE", "CONTAINERS", "CREATED", "NETWORK"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	services := []string{}
	var service map[string]interface{}
	for _, app := range apps {
		b := reqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/services", app), nil)
		check(json.Unmarshal(b, &services))
		for _, serviceID := range services {
			b := reqWant("GET", http.StatusOK, fmt.Sprintf("/applications/%s/services/%s", app, serviceID), nil)
			check(json.Unmarshal(b, &service))

			/*'%s/%s' % (application, service['name']),
			'%s@%s' % (service['repository'], service['repository_tag']),
			service['image'][:12],
			service['state'].capitalize(),
			service['container_number'],
			dateutil.parser.parse(service['creation_date']).replace(tzinfo=None),
			', '.join(ips),
			*/

			//ips := []string{}
			/* DOING PER YESNAULT for _, container := range service["containers"].(map[string]interface{}) {
				fmt.Printf("container: %+v \n", container)
				for name, network := range container["network"].(map[string]interface{}) {
					fmt.Printf("%s:%s\n", name, network["ip"])
					append(ips, fmt.Sprintf("%s:%s", name, network["ip"]))
				}
			}*/

			fmt.Fprintf(w, "%s/%s\t%s@%s\t%s\t%s\t%s\t%s\n",
				app, service["name"],
				service["repository"],
				service["repository_tag"],
				service["image"].(string)[:12],
				strings.ToUpper(service["state"].(string)),
				service["container_number"],
				service["creation_date"].(string)[:19])

			w.Flush()
		}
	}
}
