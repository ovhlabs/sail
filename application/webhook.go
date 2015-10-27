package application

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdApplicationWebhook = &cobra.Command{
	Use:   "webhook",
	Short: "Application webhook commands: sail application webhook --help",
	Long:  `Application webhook commands: sail application webhook <command>`,
}

var cmdApplicationWebhookList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the webhooks of an app: sail application webhook list <applicationName>",
	Long: `List the webhooks of an app: sail application webhook list <applicationName>
	example: sail application webhook list my-app"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application webhook list --help")
		} else {
			// Sanity
			err := internal.CheckName(args[0])
			internal.Check(err)

			internal.FormatOutputDef(internal.GetWantJSON(fmt.Sprintf("/applications/%s/hook", args[0])))
		}
	},
}

var cmdApplicationWebhookAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a webhook to an application ; sail application webhook add <applicationName> <WebhookURL>",
	Long: `Add a webhook to an application ; sail application webhook add <applicationName> <WebhookURL>
		example: sail application webhook add my-app http://www.endpoint.com/hook
		Endpoint url must except POST with json body.
		`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application webhook add --help")
		} else {
			// Sanity
			err := internal.CheckName(args[0])
			internal.Check(err)
			webhookAdd(args[0], args[1])
		}
	},
}

var cmdApplicationWebhookDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a webhook to an application ; sail application webhook delete <applicationName> <WebhookURL>",
	Long: `Delete a webhook to an application ; sail application webhook delete <applicationName> <WebhookURL>
		example: sail application webhook delete my-app http://www.endpoint.com/hook
		`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail application webhook delete --help")
		} else {
			// Sanity
			err := internal.CheckName(args[0])
			internal.Check(err)
			webhookDelete(args[0], args[1])
		}
	},
}

type webhookStruct struct {
	URL string `json:"url"`
}

func webhookAdd(namespace, webhookURL string) {

	path := fmt.Sprintf("/applications/%s/hook", namespace)

	rawBody := webhookStruct{URL: webhookURL}
	body, err := json.MarshalIndent(rawBody, " ", " ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err)
		return
	}
	internal.FormatOutputDef(internal.PostBodyWantJSON(path, body))
}

func webhookDelete(namespace, webhookURL string) {
	urlEscape := url.QueryEscape(webhookURL)

	path := fmt.Sprintf("/applications/%s/hook", namespace)
	// pass urlEscape as query string argument
	BaseURL, err := url.Parse(path)
	internal.Check(err)

	params := url.Values{}
	params.Add("url", urlEscape)

	BaseURL.RawQuery = params.Encode()
	internal.FormatOutputDef(internal.DeleteWantJSON(BaseURL.String()))
}
