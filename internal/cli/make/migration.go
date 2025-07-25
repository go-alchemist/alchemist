package make

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/go-alchemist/alchemist/internal/cli/response"
)

var MigrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Create a new SQL migration file (up and down)",
	Args:  cobra.ExactArgs(1),
	Run:   MakeMigration,
}

func init() {
	MigrationCmd.Flags().String("dir", "internal/database/migrations", "Directory for migrations")
}

func MakeMigration(cmd *cobra.Command, args []string) {
	migrationName := args[0]
	dir, _ := cmd.Flags().GetString("dir")
	if dir == "" {
		dir = "internal/database/migrations"
	}
	timestamp := time.Now().Format("20060102150405")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	upFile := fmt.Sprintf("%s/%s_%s.up.sql", dir, timestamp, migrationName)
	downFile := fmt.Sprintf("%s/%s_%s.down.sql", dir, timestamp, migrationName)

	upTemplate := "-- Write your UP SQL statements here\n"
	downTemplate := "-- Write your DOWN SQL statements here\n"

	errUp := os.WriteFile(upFile, []byte(upTemplate), 0644)
	errDown := os.WriteFile(downFile, []byte(downTemplate), 0644)

	if errUp != nil || errDown != nil {
		response.Error("Error creating migration files: " + errUp.Error() + " " + errDown.Error())
		return
	}

	response.Success(fmt.Sprintf("Migration files created:\n%s\n%s", upFile, downFile))
}
