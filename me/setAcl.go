package me

import (
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/runabove/sail/internal"
)

var cmdMeSetAcl = &cobra.Command{
	Use:   "setAcl",
	Short: "Set ip based account access restrictions: sail me setAcl <ip> [<ip> ... ]",
	Long: `Set ip based account access restrictions: sail me setAcl <ip> [<ip> ... ]
	\"example: sail me setAcl 1.2.3.4/24 4.5.6.7/32\"
	`,
	Aliases: []string{"setAcls", "set-acls", "set-acl"},
	Run: func(cmd *cobra.Command, args []string) {
		jsonStr, err := json.Marshal(args)
		internal.Check(err)
		internal.FormatOutputDef(internal.ReqWantJSON("PUT", http.StatusOK, "/user/acl", jsonStr))
	},
}
