package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sailgo/internal"
)

var cmdServiceAddLink string
var cmdServiceAddNetworkAllow string
var addPublish []string
var cmdServiceAddGateway string
var cmdServiceAddVolume string
var batch bool
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
	cmd.Flags().StringSliceVarP(&addPublish, "publish", "", nil, "Publish a container's port to the host")
	cmd.Flags().StringVarP(&cmdServiceAddGateway, "gateway", "", "", "network-input:network-output")
	cmd.Flags().StringVarP(&cmdServiceAddBody.RestartPolicy, "restart", "", "no", "{no|always[:<max>]|on-failure[:<max>]}")
	cmd.Flags().StringVarP(&cmdServiceAddVolume, "volume", "", "", "/path:size] (Size in GB)")
	cmd.Flags().BoolVarP(&batch, "batch", "", false, "do not attach console on start")
	cmd.Flags().BoolVarP(&cmdServiceAddRedeploy, "redeploy", "", false, "if the service already exists, redeploy instead")
	cmd.Flags().StringSliceVarP(&cmdServiceAddBody.ContainerEnvironment, "env", "e", nil, "override docker environment")
	//	toto = cmd.Flags().String
	return cmd
}

// PortConfig is a parameter of ServiceAdd to modify exposed container ports
type PortConfig struct {
	PublishedPort string `json:"published_port"`
	Network       string `json:"network,omitempty"`
}

// ServiceAdd struct holds all parameters sent to /applications/%s/services/%s?stream
type ServiceAdd struct {
	Service              string                       `json:"-"`
	Volumes              map[string]string            `json:"volumes"`
	Repository           string                       `json:"repository"`
	ContainerUser        string                       `json:"container_user"`
	RestartPolicy        string                       `json:"restart_policy"`
	ContainerCommand     []string                     `json:"container_command"`
	ContainerNetwork     map[string]map[string]string `json:"container_network"`
	ContainerEntrypoint  string                       `json:"container_user"`
	ContainerNumber      int                          `json:"container_number"`
	RepositoryTag        string                       `json:"repository_tag"`
	Links                map[string]map[string]string `json:"links"`
	Namespace            string                       `json:"namespace"`
	ContainerWorkdir     string                       `json:"container_workdir"`
	ContainerEnvironment []string                     `json:"container_environment"`
	ContainerModel       string                       `json:"container_model"`
	ContainerPorts       map[string][]PortConfig      `json:"container_ports"`
}

func cmdServiceAdd(cmd *cobra.Command, args []string) {
	cmdServiceAddBody.ContainerNetwork = make(map[string]map[string]string)
	cmdServiceAddBody.Links = make(map[string]map[string]string)
	cmdServiceAddBody.Volumes = make(map[string]string)
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

func serviceAdd(args ServiceAdd) {

	// Parse ContainerNetworks arguments
	for _, network := range cmdServiceAddNetwork {
		args.ContainerNetwork[network] = make(map[string]string)
	}

	// Parse ContainerPorts
	args.ContainerPorts = parsePublishedPort(addPublish)

	path := fmt.Sprintf("/applications/%s/services/%s?stream", args.Namespace, args.Service)
	body, err := json.MarshalIndent(args, " ", " ")
	if err != nil {
		fmt.Printf("Fatal: %s\n", err)
		return
	}

	if batch {
		ret := internal.ReqWant("POST", http.StatusOK, path, body)
		e := internal.DecodeError(ret)
		if e != nil {
			fmt.Printf("%s\n", e)
		} else {
			fmt.Printf("%s\n", ret)
		}
	} else {
		path = path + "?stream"
		internal.StreamWant("POST", http.StatusOK, path, body)
	}
}

func parsePublishedPort(args []string) map[string][]PortConfig {
	v := make(map[string][]PortConfig)

	for _, pub := range args {
		split := strings.Split(pub, ":")
		if len(split) == 1 { // containerPort
			v[split[0]+"/tcp"] = []PortConfig{PortConfig{PublishedPort: split[0]}}
		} else if len(split) == 2 { // network::containerPort, publishedPort:containerPort
			_, err := strconv.Atoi("-42")
			if err != nil { // network::containerPort
				key := split[0] + "/" + split[1]
				v[key] = append(v[key], PortConfig{PublishedPort: split[0], Network: split[1]})
			} else { // publishedPort:containerPort
				key := split[0] + "/tcp"
				v[key] = append(v[key], PortConfig{PublishedPort: split[1]})
			}
		} else if len(split) == 3 { // network:publishedPort:containerPort
			if split[1] == "" {
				split[1] = split[2]
			}

			key := split[1] + "/" + split[0]
			v[key] = append(v[key], PortConfig{PublishedPort: split[2], Network: split[0]})
		}
	}

	return v
}
