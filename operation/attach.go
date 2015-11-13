package operation

import (
	"fmt"
	"os"

	"github.com/runabove/sail/internal"
	"github.com/spf13/cobra"
)

var cmdOperationAttach = &cobra.Command{
	Use:   "attach",
	Short: "Attach to an ongoing operation output: sail operation attach [applicationName] <operationId>",
	Long: `Attach to an ongoing operation output: sail operation attach [applicationName] <operationId>

Example: sail operation attach devel/redis fa853ede-6c05-4823-8b20-46a5389fe0de

If the applicationName is not passed, the default application name will be used (the user's username).
`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 1:
			// applicationName was not passed. Using default one.
			applicationName := internal.GetUserName()
			operationAttach(applicationName, args[0])
		case 2:
			operationAttach(args[0], args[1])
		default:
			fmt.Fprintln(os.Stderr, "Invalid usage. sail operation attach [applicationName] <operationId>. Please see sail operation attach --help")
		}
	},
}

func operationAttach(app, operationID string) {
	// Split namespace and service
	internal.StreamPrint("GET", fmt.Sprintf("/applications/%s/operation/%s/attach", app, operationID), nil)
	internal.ExitAfterCtrlC()
}
