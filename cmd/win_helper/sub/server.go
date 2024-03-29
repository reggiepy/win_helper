package sub

import (
	"fmt"
	"github.com/spf13/cobra"
	"win_helper/pkg/server/win"
)

var (
	sId                     string
	sExecutable             string
	sName                   string
	sDescription            string
	sStartMode              string
	sDepends                string
	sLogPath                string
	sArguments              string
	sStartArguments         string
	sStopExecutable         string
	sStopArguments          string
	sEnv                    string
	sFailure                string
	sWorkingDirectory       string
	sLogMode                string
	sLogPattern             string
	sLogAutoRollAtTime      string
	sLogSizeThreshold       string
	sLogZipOlderThanNumDays string
	sLogZipDateFormat       string
)

func newServerCmd() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "generate windows exe server",
		Long:  `generate windows exe server`,
		Args:  validateServerCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := win.NewServer(
				win.WithSId(sId),
				win.WithSName(sName),
				win.WithSExecutable(sExecutable),
				win.WithSDescription(sDescription),
				win.WithSStartMode(sStartMode),
				win.WithSDepends(sDepends),
				win.WithSLogPath(sLogPath),
				win.WithSArguments(sArguments),
				win.WithSStartArguments(sStartArguments),
				win.WithSStopExecutable(sStopExecutable),
				win.WithSStopArguments(sStopArguments),
				win.WithSEnv(sEnv),
				win.WithSFailure(sFailure),
				win.WithSWorkingDirectory(sWorkingDirectory),
				win.WithSLogMode(sLogMode),
				win.WithSLogPattern(sLogPattern),
				win.WithSLogAutoRollAtTime(sLogAutoRollAtTime),
				win.WithSLogSizeThreshold(sLogSizeThreshold),
				win.WithSLogZipOlderThanNumDays(sLogZipOlderThanNumDays),
				win.WithSLogZipDateFormat(sLogZipDateFormat),
			)
			err := s.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
	serverCmd.Flags().StringVar(&sId, "id", "", "Id(default=name)")
	serverCmd.Flags().StringVar(&sName, "name", "", "name")
	serverCmd.Flags().StringVar(&sExecutable, "executable", "", "executable")
	serverCmd.Flags().StringVar(&sDescription, "description", "", "description")
	serverCmd.Flags().StringVar(&sStartMode, "startmode", "", "start mode")
	serverCmd.Flags().StringVar(&sDepends, "depends", "", "depends")
	serverCmd.Flags().StringVar(&sLogPath, "logpath", "", "log path")
	serverCmd.Flags().StringVar(&sArguments, "arguments", "", "arguments")
	serverCmd.Flags().StringVar(&sStartArguments, "startarguments", "", "start arguments")
	serverCmd.Flags().StringVar(&sStopExecutable, "stopexecutable", "", "stop executable")
	serverCmd.Flags().StringVar(&sStopArguments, "stoparguments", "", "stop arguments")
	serverCmd.Flags().StringVar(&sEnv, "env", "", "environment variables")
	serverCmd.Flags().StringVar(&sFailure, "failure", "", "failure")
	serverCmd.Flags().StringVar(&sWorkingDirectory, "workingdirectory", "", "working directory")
	serverCmd.Flags().StringVar(&sLogMode, "logmode", "", "log mode")
	serverCmd.Flags().StringVar(&sLogPattern, "logpattern", "", "log pattern")
	serverCmd.Flags().StringVar(&sLogAutoRollAtTime, "logautorollattime", "", "log auto roll at time")
	serverCmd.Flags().StringVar(&sLogSizeThreshold, "logizethreshold", "", "log size threshold")
	serverCmd.Flags().StringVar(&sLogZipOlderThanNumDays, "logzipolderthannumdays", "", "log zip older than num days")
	serverCmd.Flags().StringVar(&sLogZipDateFormat, "logzipdateformat", "", "log zip date format")

	//Boot Start ("Boot")
	//Device driver started by the operating system loader. This value is valid only for driver services.
	//System ("System")
	//Device driver started by the operating system initialization process. This value is valid only for driver services.
	//Auto Start ("Automatic")
	//Service to be started automatically by the service control manager during system startup.
	//Demand Start ("Manual")
	//Service to be started by the service control manager when a process calls the StartService method.
	//Disabled ("Disabled")
	//Service that can no longer be started.
	serverCmd.Flags().StringVar(&sStartMode, "start-mode", "", "start-mode(Boot|System|Automatic|Manual|Disabled) (default: Automatic)")

	return serverCmd
}

func validateServerCmd(cmd *cobra.Command, args []string) error {
	if sName == "" {
		return fmt.Errorf("missing name")
	}
	if sId == "" {
		sId = sName
	}
	if sExecutable == "" {
		return fmt.Errorf("missing executable")
	}
	return nil
}
