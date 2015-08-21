package me

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"stash.ovh.net/sailabove/sailgo/internal"
)

func init() {
	Cmd.AddCommand(cmdMeShow)
	Cmd.AddCommand(cmdMeSetAcl)
}

var Cmd = &cobra.Command{
	Use:     "me",
	Short:   "Me commands : sailgo me --help",
	Long:    `Me commands : sailgo me <command>`,
	Aliases: []string{"m"},
}

var cmdMeShow = &cobra.Command{
	Use:   "show",
	Short: "Show account details : sailgo me show",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(internal.GetWantJSON("/users"))
	},
}

var cmdMeSetAcl = &cobra.Command{
	Use:   "setAcl",
	Short: "Set ip based account access restrictions : sailgo me setAcl <ip> [<ip> ... ]",
	Long: `Set ip based account access restrictions : sailgo me setAcl <ip> [<ip> ... ]
	\"example : sailgo me setAcl 1.2.3.4/24 4.5.6.7/32\"
	`,
	Aliases: []string{"setAcls", "set-acls", "set-acl"},
	Run: func(cmd *cobra.Command, args []string) {
		jsonStr, err := json.Marshal(args)
		internal.Check(err)
		fmt.Println(internal.ReqWantJSON("PUT", http.StatusOK, "/user/acl", jsonStr))
	},
}
