package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"win_helper/pkg/server/win"
)

var (
	name       string
	executable string
	onfailure  string
	logMode    string
)

func newServerCmd() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "generate windows exe server",
		Long:  `generate windows exe server`,
		Args:  validateServerCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := win.NewServer(
				win.WithName(name),
				win.WithBasePath(baseDir),
				win.WithExecutable(executable),
				win.WithLogMode(logMode),
			)
			err := s.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
	serverCmd.Flags().StringVar(&name, "name", "", "name")
	serverCmd.Flags().StringVar(&executable, "executable", "", "executable")
	serverCmd.Flags().StringVar(&onfailure, "onfailure", "restart", "onfailure")
	serverCmd.Flags().StringVar(&logMode, "log-mode", "roll", "log-mode")

	return serverCmd
}

func validateServerCmd(cmd *cobra.Command, args []string) error {
	if name == "" {
		return fmt.Errorf("missing name")
	}
	if executable == "" {
		return fmt.Errorf("missing executable")
	}
	return nil
}
