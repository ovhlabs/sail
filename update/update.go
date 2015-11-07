package update

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

// used by CI to inject architecture (linux-amd64, etc...) at build time
var architecture string
var urlGitubReleases = "https://github.com/runabove/sail/releases"

// Cmd update
var Cmd = &cobra.Command{
	Use:     "update",
	Short:   "Update sail to the latest release version: sail update",
	Long:    `sail update`,
	Aliases: []string{"up"},
	Run: func(cmd *cobra.Command, args []string) {
		doUpdate("", architecture)
	},
}

func getURLArtifactFromGithub(architecture string) string {
	client := github.NewClient(nil)
	release, resp, err := client.Repositories.GetLatestRelease("runabove", "sail")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Repositories.GetLatestRelease returned error: %v\n%v", err, resp.Body)
		os.Exit(1)
	}

	if len(release.Assets) > 0 {
		for _, asset := range release.Assets {
			if *asset.Name == "sail-"+architecture {
				return *asset.BrowserDownloadURL
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Invalid Artifacts on latest release. Please try again in few minutes.\n")
	fmt.Fprintf(os.Stderr, "If the problem persists, please open an issue on https://github.com/runabove/sail/issues\n")
	os.Exit(1)
	return ""
}

func getContentType(resp *http.Response) string {
	for k, v := range resp.Header {
		if k == "Content-Type" && len(v) >= 1 {
			return v[0]
		}
	}
	return ""
}

func doUpdate(baseurl, architecture string) {
	if architecture == "" {
		fmt.Fprintf(os.Stderr, "You seem to have a custom build of sail\n")
		fmt.Fprintf(os.Stderr, "Please download latest release on %s\n", urlGitubReleases)
		os.Exit(1)
	}

	url := getURLArtifactFromGithub(architecture)
	if internal.Verbose {
		fmt.Printf("Url to update sail: %s\n", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when downloading sail: %s\n", err.Error())
		fmt.Printf("Url: %s\n", url)
		os.Exit(1)
	}

	contentType := getContentType(resp)
	if contentType != "application/octet-stream" {
		fmt.Fprintf(os.Stderr, "Invalid Binary (Content-Type: %s). Please try again or download it manually from %s\n", contentType, urlGitubReleases)
		fmt.Printf("Url: %s\n", url)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Error http code: %d, url called: %s\n", resp.StatusCode, url)
		os.Exit(1)
	}

	fmt.Printf("Getting latest release from : %s ...\n", url)
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when updating sail: %s\n", err.Error())
		fmt.Fprintf(os.Stderr, "Url: %s\n", url)
		os.Exit(1)
	}
	fmt.Println("Update done.")
}
