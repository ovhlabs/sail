package me

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdMeSSHKey = &cobra.Command{
	Use:   "sshkey",
	Short: "User ssh key commands: sail me sshkey --help",
	Long:  `User ssh key commands: sail me sshkey <command>`,
}

var cmdMeSSHKeyList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the ssh keys of this account: sail me sshkey list",
	Long:    `List the ssh keys of this account: sail me sshkey list`,
	Run: func(cmd *cobra.Command, args []string) {
		sshKeyList()
	},
}

var cmdMeSSHKeyAdd = &cobra.Command{

	Use:     "add",
	Aliases: []string{"register"},
	Short:   "Add an ssh key this account from a public key file or from stdin: sail me sshkey add <keyName> [<filepath>] ",
	Long: `Add an ssh key this account from a public key file or from stdin: sail me sshkey add <keyName> [<filepath>]
Support spaces in key name.
examples:
sail me sshkey add "my test name" /home/user/.ssh/id_rsa
sail me sshkey add "my test name" then copy and paste the content of /home/user/.ssh/id_rsa
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail me sshkey add --help")
		} else if len(args) == 1 {
			//a name without a filepath, try stdin
			fmt.Println("Please paste the content of your key file or press Ctrl+c to cancel.")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}
			keyLine := scanner.Text()
			sshKeyAdd(keyLine, args[0])
		} else {
			// a name and a filepath, read the file and go on
			keyLine, err := ioutil.ReadFile(args[1])
			internal.Check(err)

			sshKeyAdd(string(keyLine), args[0])
		}
	},
}

var cmdMeSSHKeyDelete = &cobra.Command{

	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "Delete an ssh keys from the account: sail me sshkey delete <fingerprint>",
	Long: `Delete an ssh keys of this account: sail me sshkey delete <fingerprint>
	example : sail me sshkey delete 0d/keDZZb3OAjj+8JI7T5iIxMLUT643YfW3mBznqrC8=`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "" {
			fmt.Fprintln(os.Stderr, "Invalid usage. Please see sail me sshkey delete --help")
		} else {
			sshKeyDelete(args[0])
		}
	},
}

type sshkeyStruct struct {
	KeyLine string `json:"key_line"`
	KeyName string `json:"key_name"`
}

func sshKeyList() {

	type keyStruct struct {
		Name        string `json:"name"`
		Fingerprint string `json:"fingerprint"`
		PublicKey   string `json:"public_key"`
	}

	b := internal.ReqWant("GET", http.StatusOK, "/user/keys", nil)
	var keys []keyStruct

	err := json.Unmarshal(b, &keys)
	internal.Check(err)

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	titles := []string{"NAME", "FINGERPRINT"}
	fmt.Fprintln(w, strings.Join(titles, "\t"))

	for _, key := range keys {
		fmt.Fprintf(w, "%s\t%s\n",
			key.Name,
			key.Fingerprint,
		)
	}
	w.Flush()
}

func sshKeyAdd(keyLine, keyName string) {
	path := "/user/keys"

	rawBody := sshkeyStruct{
		KeyLine: string(keyLine),
		KeyName: string(keyName),
	}
	body, err := json.MarshalIndent(rawBody, " ", " ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err)
		return
	}

	internal.FormatOutputDef(internal.PostBodyWantJSON(path, body))
}

func sshKeyDelete(fingerprint string) {
	urlEscape := url.QueryEscape(fingerprint)

	path := "/user/keys"

	// pass urlEscape as query string argument
	BaseURL, err := url.Parse(path)
	internal.Check(err)

	params := url.Values{}
	params.Add("fingerprint", urlEscape)

	BaseURL.RawQuery = params.Encode()
	internal.FormatOutputDef(internal.DeleteWantJSON(BaseURL.String()))
}
