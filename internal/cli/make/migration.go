package make

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"

	"github.com/go-alchemist/alchemist/internal/cli/config"
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

	flagDir, _ := cmd.Flags().GetString("dir")
	originalDir := flagDir
	if originalDir == "" {
		originalDir = config.Config.GetString("paths.migrations")
		if originalDir == "" {
			originalDir = "internal/database/migrations"
		}
	}

	dir := response.SelectMicroserviceIfEnabled()
	if dir == "" {
		response.Error("Microservice feature is not enabled or directory not found")
		return
	}

	dir = path.Join(dir, originalDir)

	timestamp := time.Now().Format("20060102150405")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	upFile := path.Join(dir, fmt.Sprintf("%s_%s.up.sql", timestamp, migrationName))
	downFile := path.Join(dir, fmt.Sprintf("%s_%s.down.sql", timestamp, migrationName))

	upTemplate := "-- Write your UP SQL statements here\n"
	downTemplate := "-- Write your DOWN SQL statements here\n"

	errUp := os.WriteFile(upFile, []byte(upTemplate), 0644)
	errDown := os.WriteFile(downFile, []byte(downTemplate), 0644)

	if errUp != nil || errDown != nil {
		response.Error(fmt.Sprintf("Error creating migration files: %v %v", errUp, errDown))
		return
	}

	response.Success(fmt.Sprintf("Migration files created:\n%s\n%s", upFile, downFile))
}
