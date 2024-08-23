package sub

import (
	"fmt"
	"os"

	"github.com/jarvanstack/mysqldump"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mysqlCmd)
	mysqlCmd.AddCommand(dumpMysqlCmd)
	mysqlCmd.AddCommand(sourceMysqlCmd)

	dumpMysqlCmd.Flags().String("dsn", "", "Database dsn")
	dumpMysqlCmd.Flags().String("dest", "", "dest sql")

	sourceMysqlCmd.Flags().String("dsn", "", "Database dsn")
	sourceMysqlCmd.Flags().String("dest", "", "dest sql")
}

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql db command",
	Long:  "mysql db command",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.UsageString())
		return nil
	},
}

var dumpMysqlCmd = &cobra.Command{
	Use:   "dump",
	Short: "dump mysql db command",
	Long:  "dump mysql db command",
	Args: func(cmd *cobra.Command, args []string) error {
		dns, _ := cmd.Flags().GetString("dsn")
		if dns == "" {
			return fmt.Errorf("missing dsn")
		}
		dest, _ := cmd.Flags().GetString("dest")
		if dest == "" {
			return fmt.Errorf("missing dest")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dest, _ := cmd.Flags().GetString("dest")
		dns, _ := cmd.Flags().GetString("dsn")
		f, _ := os.Create(dest)

		_ = mysqldump.Dump(
			dns,                       // DSN
			mysqldump.WithDropTable(), // Option: Delete table before create (Default: Not delete table)
			mysqldump.WithData(),      // Option: Dump Data (Default: Only dump table schema)
			// mysqldump.WithTables("test"), // Option: Dump Tables (Default: All tables)
			mysqldump.WithWriter(f), // Option: Writer (Default: os.Stdout)
		)
		return nil
	},
}

var sourceMysqlCmd = &cobra.Command{
	Use:   "source",
	Short: "source mysql db command",
	Long:  "source mysql db command",
	Args: func(cmd *cobra.Command, args []string) error {
		dns, _ := cmd.Flags().GetString("dsn")
		if dns == "" {
			return fmt.Errorf("missing dsn")
		}
		dest, _ := cmd.Flags().GetString("dest")
		if dest == "" {
			return fmt.Errorf("missing dest")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dest, _ := cmd.Flags().GetString("dest")
		dns, _ := cmd.Flags().GetString("dsn")
		f, _ := os.Open(dest)
		_ = mysqldump.Source(
			dns, // DSN
			f,
			mysqldump.WithMergeInsert(1000), // Option: Merge insert 1000 (Default: Not merge insert)
			mysqldump.WithDebug(),           // Option: Print execute sql (Default: Not print execute sql)
		)
		return nil
	},
}
