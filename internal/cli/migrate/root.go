package migrate

import (
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  "Manage and run SQL database migrations.",
}

func init() {
	MigrateCmd.AddCommand(RunCmd)
}
