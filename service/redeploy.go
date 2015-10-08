package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/google/shlex"

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
	redeployPool         string
	redeployCommand      string
	redeployEntrypoint   string
)

func redeployCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "redeploy",
		Short:   "Redeploy a docker service: sail service redeploy <applicationName>/<serviceId>",
		Long:    `Redeploy a docker service: sail service redeploy <applicationName>/<serviceId>`,
		Aliases: []string{"restart"},
		Run:     cmdRedeploy,
	}

	cmd.Flags().StringVarP(&redeployBody.ContainerModel, "model", "", "", "Container model")
	cmd.Flags().IntVarP(&redeployBody.ContainerNumber, "number", "", 0, "Number of container to run")
	cmd.Flags().StringSliceVarP(&redeployBody.ContainerEnvironment, "env", "e", nil, "override docker environment")
	cmd.Flags().StringVarP(&redeployBody.RestartPolicy, "restart", "", "", "{no|always[:<max>]|on-failure[:<max>]}")
	cmd.Flags().StringVarP(&redeployCommand, "command", "", "", "override docker run command")
	cmd.Flags().StringVarP(&redeployBody.RepositoryTag, "tag", "", "", "deploy from new image version")
	cmd.Flags().StringVarP(&redeployBody.ContainerWorkdir, "workdir", "", "", "override docker workdir")
	cmd.Flags().StringVarP(&redeployEntrypoint, "entrypoint", "", "", "override docker entrypoint")
	cmd.Flags().StringVarP(&redeployBody.ContainerUser, "user", "", "", "override docker user")
	cmd.Flags().StringSliceVar(&redeployNetwork, "network", nil, "public|private|<namespace name>")
	cmd.Flags().StringVarP(&redeployNetworkAllow, "network-allow", "", "", "[network:]ip[/mask] Use IPs whitelist")
	cmd.Flags().StringSliceVarP(&redeployPublished, "publish", "", nil, "Publish a container's port to the host")
	cmd.Flags().StringSliceVarP(&redeployGateway, "gateway", "", nil, "network-input:network-output")
	cmd.Flags().StringSliceVarP(&redeployVolume, "volume", "", nil, "/path:size] (Size in GB)")
	cmd.Flags().BoolVarP(&redeployBatch, "batch", "", false, "do not attach console on start")
	cmd.Flags().StringSliceVarP(&redeployLink, "link", "", nil, "name:alias")
	cmd.Flags().StringVarP(&redeployPool, "pool", "", "", "Dedicated host pool")
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
	ContainerEntrypoint  []string                       `json:"container_entrypoint,omitempty"`
	ContainerNumber      int                            `json:"container_number,omitempty"`
	RepositoryTag        string                         `json:"repository_tag,omitempty"`
	Links                map[string]string              `json:"links,omitempty"`
	Application          string                         `json:"namespace,omitempty"`
	ContainerWorkdir     string                         `json:"container_workdir,omitempty"`
	ContainerEnvironment []string                       `json:"container_environment,omitempty"`
	ContainerModel       string                         `json:"container_model,omitempty"`
	ContainerPorts       map[string][]PortConfig        `json:"container_ports,omitempty"`
	Pool                 string                         `json:"pool,omitempty"`
}

func cmdRedeploy(cmd *cobra.Command, args []string) {
	usage := "Invalid usage. sail service redeploy <applicationName>/<serviceId>. Please see sail service redeploy --help\n"
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, usage)
		return
	}

	// Split namespace and repository
	host, app, service, _, err := internal.ParseResourceName(args[0])
	internal.Check(err)
	redeployBody.Application = app
	redeployBody.Service = service

	if !internal.CheckHostConsistent(host) {
		fmt.Fprintf(os.Stderr, "Error: Invalid Host %s for endpoint %s\n", host, internal.Host)
		os.Exit(1)
	}

	// Redeploy
	serviceRedeploy(redeployBody)
}

func serviceRedeploy(args Redeploy) {

	// Parse volumes
	if len(redeployVolume) > 0 {
		args.Volumes = make(map[string]VolumeConfig)
	}

	// Parse command
	if redeployCommand != "" {
		command, err := shlex.Split(redeployCommand)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal, cannot split command %s\n", err)
			return
		}
		args.ContainerCommand = command
	}

	// Parse Entrypoint
	if redeployEntrypoint != "" {
		entrypoint, err := shlex.Split(redeployEntrypoint)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal, cannot split command %s\n", err)
			return
		}
		args.ContainerEntrypoint = entrypoint
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
	if len(redeployLink) > 0 {
		args.Links = make(map[string]string)
	}

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
			fmt.Fprintf(os.Stderr, "Invalid gateway parameter, should be \"input:output\". Typically, output will be 'predictor' or 'public'")
			os.Exit(1)
		}
	}
	// Load Pool
	args.Pool = redeployPool

	// Parse ContainerPorts
	args.ContainerPorts = parsePublishedPort(redeployPublished)
	app := args.Application
	service := args.Service

	path := fmt.Sprintf("/applications/%s/services/%s/redeploy", app, service)
	body, err := json.MarshalIndent(args, " ", " ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err)
		return
	}

	// Attach console
	if !redeployBatch {
		internal.StreamPrint("GET", fmt.Sprintf("/applications/%s/services/%s/attach", app, service), nil)
	}

	// Redeploy
	buffer, _, err := internal.Stream("POST", path+"?stream", body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	line, err := internal.DisplayStream(buffer)
	internal.Check(err)
	if line != nil {
		var data map[string]interface{}
		err = json.Unmarshal(line, &data)
		internal.Check(err)

		fmt.Printf("Hostname: %v\n", data["hostname"])
		fmt.Printf("Running containers: %v/%v\n", data["container_number"], data["container_target"])
	}

	if !redeployBatch {
		internal.ExitAfterCtrlC()
	}
}
