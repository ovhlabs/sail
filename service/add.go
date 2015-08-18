package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"stash.ovh.net/sailabove/sailgo/Godeps/_workspace/src/github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sailgo/internal"
)

var cmdServiceAddLink string
var cmdServiceAddNetworkAllow string
var cmdServiceAddPublish string
var cmdServiceAddGateway string
var cmdServiceAddVolume string
var cmdServiceAddBatch bool
var cmdServiceAddRedeploy bool
var cmdServiceAddBody ServiceAdd
var cmdServiceAddNetwork []string

func addCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new docker service",
		Long: `add [<namespace>/]<repository>[:tag] [namespace/]<service-name>
		--model         Container model
		--number        Number of container to run
		[--link         name:alias]
		[--network      {public|private|<namespace name>}]
		[--network-allow [network:]ip[/mask] Use IPs whitelist]
		[--publish, -p  Publish a container's port to the host]
		[                 format: network:publishedPort:containerPort, network::containerPort, publishedPort:containerPort, containerPort]
		[--gateway      network-input:network-output
		[--restart {no|always[:<max>]|on-failure[:<max>]}]
		[--volume       /path:size] (Size in GB)
		[--batch        do not attach console on start]
		[--redeploy     if the service already exists, redeploy instead]
		
		override docker options:
			--user
			--entrypoint
			--command
			--workdir
			--environment KEY=val
		other options:
		`,
		Run: cmdServiceAdd,
	}
	cmd.Flags().StringVarP(&cmdServiceAddBody.ContainerModel, "model", "", "x1", "Container model")
	cmd.Flags().IntVarP(&cmdServiceAddBody.ContainerNumber, "number", "", 1, "Number of container to run")
	cmd.Flags().StringVarP(&cmdServiceAddLink, "link", "", "", "name:alias")
	cmd.Flags().StringSliceVar(&cmdServiceAddNetwork, "network", []string{"public", "private"}, "public|private|<namespace name>")
	cmd.Flags().StringVarP(&cmdServiceAddNetworkAllow, "network-allow", "", "", "[network:]ip[/mask] Use IPs whitelist")
	cmd.Flags().StringVarP(&cmdServiceAddPublish, "publish", "P", "", "Publish a container's port to the host")
	cmd.Flags().StringVarP(&cmdServiceAddGateway, "gateway", "", "", "network-input:network-output")
	cmd.Flags().StringVarP(&cmdServiceAddBody.RestartPolicy, "restart", "", "no", "{no|always[:<max>]|on-failure[:<max>]}")
	cmd.Flags().StringVarP(&cmdServiceAddVolume, "volume", "", "", "/path:size] (Size in GB)")
	cmd.Flags().BoolVarP(&cmdServiceAddBatch, "batch", "", false, "do not attach console on start")
	cmd.Flags().BoolVarP(&cmdServiceAddRedeploy, "redeploy", "", false, "if the service already exists, redeploy instead")

	return cmd
}

// PortConfig is a parameter of ServiceAdd to modify exposed container ports
type PortConfig struct {
	PublishedPort string `json:"published_port"`
}

// ServiceAdd struct holds all parameters sent to /applications/%s/services/%s?stream
type ServiceAdd struct {
	Service              string                       `json:"-"`
	Namespace            string                       `json:"namespace"`
	Repository           string                       `json:"repository"`
	RepositoryTag        string                       `json:"repository_tag"`
	ContainerModel       string                       `json:"container_model"`
	ContainerNumber      int                          `json:"container_number"`
	ContainerUser        string                       `json:"container_user"`
	ContainerEntrypoint  string                       `json:"container_user"`
	ContainerCommand     []string                     `json:"container_command"`
	ContainerWorkdir     string                       `json:"container_workdir"`
	ContainerEnvironment []string                     `json:"container_environment"`
	ContainerNetwork     map[string]map[string]string `json:"container_network"`
	ContainerPorts       map[string][]PortConfig      `json:"container_ports"`
	Links                map[string]map[string]string `json:"links"`
	Volumes              map[string]string            `json:"volumes"`
	RestartPolicy        string                       `json:"restart_policy"`
}

func cmdServiceAdd(cmd *cobra.Command, args []string) {
	cmdServiceAddBody.ContainerNetwork = make(map[string]map[string]string)
	cmdServiceAddBody.Links = make(map[string]map[string]string)
	cmdServiceAddBody.Volumes = make(map[string]string)
	cmdServiceAddBody.ContainerEnvironment = make([]string, 0)
	cmdServiceAddBody.ContainerPorts = make(map[string][]PortConfig)

	if len(args) != 2 {
		fmt.Printf("Invalid usage. sailgo service add <application>/<repository>[tag] <service>. Please see sailgo service add --help\n")
		return
	}

	// Get args
	cmdServiceAddBody.Repository = args[0]
	cmdServiceAddBody.Service = args[1]

	// Split repo URL and tag
	split := strings.Split(cmdServiceAddBody.Repository, ":")
	if len(split) > 1 {
		cmdServiceAddBody.Repository = split[0]
		cmdServiceAddBody.RepositoryTag = split[1]
	}

	// Split namespace and repository
	split = strings.Split(cmdServiceAddBody.Repository, "/")
	if len(split) > 1 {
		cmdServiceAddBody.Namespace = split[0]
		cmdServiceAddBody.Repository = split[1]
	}

	serviceAdd(cmdServiceAddBody)
}

/* curl -f -s -XPOST -HContent-Type:application/json -H "Authorization: Basic ${bamboo.auth_base64_password}" -
d '{
				"volumes": null,
				"repository": "api",
				"container_user": null,
				"restart_policy": "always",
				"container_command": null,
				"container_network": {"predictor": {}},
				"container_entrypoint": null,
				"container_number": 1,
				"repository_tag": "'"${branch}"'",
				"links": {},
				"namespace": "apiorder",
				"container_workdir": null,
				"container_environment": [
					"DEBUG=${bamboo.API_DEBUG}",
					"SECRET_KEY=${bamboo.API_SECRET_KEY_PASSWORD}",
					"SQL_NAME=${bamboo.API_SQL_NAME}",
					"SQL_USER=${bamboo.API_SQL_USER}",
					"IS_IN_DOCKER=${bamboo.IS_IN_DOCKER}",
					"SQL_PASS=${bamboo.API_SQL_PASS_PASSWORD}",
					"SQL_HOST=${bamboo.API_SQL_HOST}",
					"SQLALCHEMY_DATABASE_URI=${bamboo.API_SQLALCHEMY_DATABASE_URI}",
					"SHARED_KEY_API_RIP=${bamboo.API_SHARED_KEY_API_RIP_PASSWORD}",
					"GRAYLOG_HOST=${bamboo.GRAYLOG_HOST}",
					"GRAYLOG_PORT=${bamboo.GRAYLOG_PORT}",
					"GRAYLOG_TLS=${bamboo.GRAYLOG_TLS}",
					"GRAYLOG_FLAG=${bamboo.API_GRAYLOG_FLAG}",
					"ALERT_ADMINS=${bamboo.API_ALERT_ADMINS}",
					"ALERT_FROM=${bamboo.API_ALERT_FROM}",
					"SMTP_HOST=${bamboo.SMTP_HOST}",
					"SMTP_USER=${bamboo.SMTP_USER}",
					"SMTP_PASS=${bamboo.SMTP_PASS_PASSWORD}",
					"SMTP_PORT=${bamboo.SMTP_PORT}"
				],
				"container_model": "x1",
				"container_ports": {
								"5000/tcp": [
									{
													"published_port": "5000"
									}]
					}
	}'
https://p19-1.sailabove.io/v1/applications/apiorder/services/api
*/
func serviceAdd(args ServiceAdd) {

	// Parse ContainerNetworks arguments
	for _, network := range cmdServiceAddNetwork {
		args.ContainerNetwork[network] = make(map[string]string)
	}

	// Parse ContainerPorts
	args.ContainerPorts["80/tcp"] = []PortConfig{PortConfig{PublishedPort: "80"}}

	path := fmt.Sprintf("/applications/%s/services/%s?stream", args.Namespace, args.Service)
	body, err := json.Marshal(args)
	if err != nil {
		fmt.Printf("Fatal: %s\n", err)
		return
	}
	internal.StreamWant("POST", http.StatusOK, path, body)
}
