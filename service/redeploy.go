package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/runabove/sail/internal"
)

var (
	redeployBody         Redeploy
	redeployPublished    []string
	redeployLink         []string
	redeployNetwork      []string
	redeployNetworkAllow string
	redeployGateway      []string
	redeployVolume       []string
	redeployBatch        bool
)

func redeployCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "redeploy",
		Short: "Redeploy a docker service: sail service redeploy <applicationName>/<serviceId>",
		Long:  `Redeploy a docker service: sail service redeploy <applicationName>/<serviceId>`,
		Run:   cmdRedeploy,
	}

	cmd.Flags().StringVarP(&redeployBody.ContainerModel, "model", "", "", "Container model")
	cmd.Flags().IntVarP(&redeployBody.ContainerNumber, "number", "", 0, "Number of container to run")
	cmd.Flags().StringSliceVarP(&redeployBody.ContainerEnvironment, "env", "e", nil, "override docker environment")
	cmd.Flags().StringVarP(&redeployBody.RestartPolicy, "restart", "", "", "{no|always[:<max>]|on-failure[:<max>]}")
	cmd.Flags().StringSliceVarP(&redeployBody.ContainerCommand, "command", "", nil, "override docker run command")
	cmd.Flags().StringVarP(&redeployBody.RepositoryTag, "tag", "", "", "deploy from new image version")
	cmd.Flags().StringVarP(&redeployBody.ContainerWorkdir, "workdir", "", "", "override docker workdir")
	cmd.Flags().StringVarP(&redeployBody.ContainerEntrypoint, "entrypoint", "", "", "override docker entrypoint")
	cmd.Flags().StringVarP(&redeployBody.ContainerUser, "user", "", "", "override docker user")

	cmd.Flags().StringSliceVar(&redeployNetwork, "network", nil, "public|private|<namespace name>")
	cmd.Flags().StringVarP(&redeployNetworkAllow, "network-allow", "", "", "[network:]ip[/mask] Use IPs whitelist")
	cmd.Flags().StringSliceVarP(&redeployPublished, "publish", "", nil, "Publish a container's port to the host")
	cmd.Flags().StringSliceVarP(&redeployGateway, "gateway", "", nil, "network-input:network-output")
	cmd.Flags().StringSliceVarP(&redeployVolume, "volume", "", nil, "/path:size] (Size in GB)")
	cmd.Flags().BoolVarP(&redeployBatch, "batch", "", false, "do not attach console on start")
	cmd.Flags().StringSliceVarP(&redeployLink, "link", "", nil, "name:alias")
	return cmd
}

// Redeploy struct holds all parameters sent to /applications/%s/services/%s/redeploy
type Redeploy struct {
	Service              string                         `json:"-"`
	Volumes              map[string]VolumeConfig        `json:"volumes,omitempty"`
	Repository           string                         `json:"repository,omitempty"`
	ContainerUser        string                         `json:"container_user,omitempty"`
	RestartPolicy        string                         `json:"restart_policy,omitempty"`
	ContainerCommand     []string                       `json:"container_command,omitempty"`
	ContainerNetwork     map[string]map[string][]string `json:"container_network,omitempty"`
	ContainerEntrypoint  string                         `json:"container_user,omitempty"`
	ContainerNumber      int                            `json:"container_number,omitempty"`
	RepositoryTag        string                         `json:"repository_tag,omitempty"`
	Links                map[string]string              `json:"links,omitempty"`
	Application          string                         `json:"namespace,omitempty"`
	ContainerWorkdir     string                         `json:"container_workdir,omitempty"`
	ContainerEnvironment []string                       `json:"container_environment,omitempty"`
	ContainerModel       string                         `json:"container_model,omitempty"`
	ContainerPorts       map[string][]PortConfig        `json:"container_ports,omitempty"`
}

func cmdRedeploy(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail service redeploy <applicationName>/<serviceId>. Please see sail service redeploy --help\n"
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	split := strings.Split(args[0], "/")
	if len(split) != 2 {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	// Get args
	redeployBody.Application = split[0]
	redeployBody.Service = split[1]
	serviceRedeploy(redeployBody)
}

func serviceRedeploy(args Redeploy) {

	// Parse volumes
	if len(redeployVolume) > 0 {
		args.Volumes = make(map[string]VolumeConfig)
	}

	for _, vol := range redeployVolume {
		t := strings.Split(vol, ":")
		if len(t) == 2 {
			args.Volumes[t[0]] = VolumeConfig{Size: t[1]}
		} else if len(t) == 1 {
			args.Volumes[t[0]] = VolumeConfig{Size: "10"}
		} else {
			fmt.Fprintf(os.Stderr, "Error: Volume parameter '%s' not formated correctly\n", vol)
			os.Exit(1)
		}
	}

	// Parse links
	for _, link := range redeployLink {
		t := strings.Split(link, ":")
		if len(t) == 1 {
			args.Links[t[0]] = t[0]
		} else {
			args.Links[t[0]] = t[1]
		}
	}

	// Parse ContainerNetworks arguments
	if len(redeployNetwork) > 0 {
		args.ContainerNetwork = make(map[string]map[string][]string)
	}
	for _, network := range redeployNetwork {
		args.ContainerNetwork[network] = make(map[string][]string)
	}

	for _, gat := range redeployGateway {
		t := strings.Split(gat, ":")
		if len(t) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid gateway parameter, should be \"input:output\"")
			os.Exit(1)
		}
		if _, ok := args.ContainerNetwork[t[0]]; !ok {
			fmt.Fprintf(os.Stderr, "Not configured input network %s\n", t[0])
			os.Exit(1)
		}
		if _, ok := args.ContainerNetwork[t[1]]; !ok {
			fmt.Fprintf(os.Stderr, "Not configured onput network %s\n", t[1])
			os.Exit(1)
		}
		args.ContainerNetwork[t[0]]["gateway_to"] = append(args.ContainerNetwork[t[0]]["gateway_to"], t[1])
	}

	// Parse ContainerPorts
	args.ContainerPorts = parsePublishedPort(redeployPublished)

	path := fmt.Sprintf("/applications/%s/services/%s/redeploy", args.Application, args.Service)
	body, err := json.MarshalIndent(args, " ", " ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err)
		return
	}

	buffer, _, err := internal.Stream("POST", path+"?stream", body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	err = internal.DisplayStream(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
