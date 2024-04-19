package sub

import (
	"fmt"
	"github.com/jarvanstack/mysqldump"
	"github.com/spf13/cobra"
	"os"
)

var (
	MY_Dsn    string = ""
	MY_Source string = ""
	MY_Dest   string = ""
)

func newMysqlCmd() *cobra.Command {
	mysqlCmd := &cobra.Command{
		Use:   "mysql",
		Short: "mysql db command",
		Long:  "mysql db command",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(cmd.UsageString())
			return nil
		},
	}
	mysqlCmd.AddCommand(newDumpMysqlCmd())
	mysqlCmd.AddCommand(newSourceMysqlCmd())
	return mysqlCmd
}

func newDumpMysqlCmd() *cobra.Command {
	dumpMysqlCmd := &cobra.Command{
		Use:   "dump",
		Short: "dump mysql db command",
		Long:  "dump mysql db command",
		Args:  ValidateDumpMysqlCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			f, _ := os.Create(MY_Dest)
			_ = mysqldump.Dump(
				MY_Dsn,                    // DSN
				mysqldump.WithDropTable(), // Option: Delete table before create (Default: Not delete table)
				mysqldump.WithData(),      // Option: Dump Data (Default: Only dump table schema)
				//mysqldump.WithTables("test"), // Option: Dump Tables (Default: All tables)
				mysqldump.WithWriter(f), // Option: Writer (Default: os.Stdout)
			)
			return nil
		},
	}
	dumpMysqlCmd.Flags().StringVar(&MY_Dsn, "dsn", "", "Database dsn")
	dumpMysqlCmd.Flags().StringVar(&MY_Dest, "dest", "", "dest sql")
	return dumpMysqlCmd
}
func ValidateDumpMysqlCmd(cmd *cobra.Command, args []string) error {
	if MY_Dsn == "" {
		return fmt.Errorf("missing dsn")
	}
	if MY_Dest == "" {
		return fmt.Errorf("missing dest")
	}
	return nil
}
func newSourceMysqlCmd() *cobra.Command {
	sourceMysqlCmd := &cobra.Command{
		Use:   "source",
		Short: "source mysql db command",
		Long:  "source mysql db command",
		Args:  ValidateSourceMysqlCmd,
		RunE: func(cmd *cobra.Command, args []string) error {
			f, _ := os.Open(MY_Source)
			_ = mysqldump.Source(
				MY_Dsn, // DSN
				f,
				mysqldump.WithMergeInsert(1000), // Option: Merge insert 1000 (Default: Not merge insert)
				mysqldump.WithDebug(),           // Option: Print execute sql (Default: Not print execute sql)
			)
			return nil
		},
	}
	sourceMysqlCmd.Flags().StringVar(&MY_Dsn, "dsn", "", "Database dsn")
	sourceMysqlCmd.Flags().StringVar(&MY_Source, "dest", "", "dest sql")
	return sourceMysqlCmd
}
func ValidateSourceMysqlCmd(cmd *cobra.Command, args []string) error {
	if MY_Dsn == "" {
		return fmt.Errorf("missing dsn")
	}
	if MY_Source == "" {
		return fmt.Errorf("missing dest")
	}
	return nil
}
