package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"win_helper/pkg/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions. `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("current version v%s -- HEAD", version.Full())
	},
}
