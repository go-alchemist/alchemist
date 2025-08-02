package make

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/orochaa/go-clack/prompts"
	"github.com/urfave/cli/v2"

	"github.com/go-alchemist/alchemist/internal/cli/utils"
)

func MakeMigration(c *cli.Context) error {
	migrationName, err := getNamePrompt(c, "Migration")
	if err != nil {
		return err
	}

	base := "."
	structure := utils.DetectProjectStructure(base)
	serviceDir := utils.SelectMicroserviceIfEnabled(structure)
	if serviceDir == "" {
		return utils.PrintErrorAndReturn("Microservice not found or not enabled")
	}

	sDir := c.String("dir")
	targetDir, err := utils.GetTargetDir(base, structure, serviceDir, sDir, MigrationTargetDir)
	if err != nil {
		return utils.PrintErrorAndReturn("Could not determine target directory for migration")
	}

	if err := utils.EnsureDir(targetDir); err != nil {
		return utils.PrintErrorAndReturn("Could not create the migration directory. Please check your permissions and try again.")
	}

	timestamp := time.Now().Format("20060102150405")
	safeName := strings.ReplaceAll(strings.ToLower(migrationName), " ", "_")
	upFile := filepath.Join(targetDir, fmt.Sprintf("%s_%s.up.sql", timestamp, safeName))
	downFile := filepath.Join(targetDir, fmt.Sprintf("%s_%s.down.sql", timestamp, safeName))

	if utils.FileExists(upFile) || utils.FileExists(downFile) {
		return utils.PrintErrorAndReturn("A migration with this name already exists. Choose another name.")
	}

	upTemplate := "-- Write your UP SQL statements here\n"
	downTemplate := "-- Write your DOWN SQL statements here\n"

	if err := os.WriteFile(upFile, []byte(upTemplate), 0644); err != nil {
		return utils.PrintErrorAndReturn(fmt.Sprintf("Error creating up migration file: %v", err))
	}
	if err := os.WriteFile(downFile, []byte(downTemplate), 0644); err != nil {
		return utils.PrintErrorAndReturn(fmt.Sprintf("Error creating down migration file: %v", err))
	}

	prompts.Success(fmt.Sprintf("Migration %s created successfully at:", migrationName))
	utils.PrintSuccessf("\n  Up Migration: %s\n  Down Migration: %s", upFile, downFile)
	return nil
}

func MigrationTargetDir(base, structure, service string) (string, error) {
	switch structure {
	case "modular":
		return filepath.Join(base, service, "database", "migrations"), nil
	case "domain_driven":
		return filepath.Join(base, service, "database", "migrations"), nil
	case "clean_architecture":
		return filepath.Join(base, service, "infrastructure", "migrations"), nil
	case "layered":
		return filepath.Join(base, service, "internal", "database", "migrations"), nil
	case "custom":
		return utils.SelectCustomDirectory(filepath.Join(base, service), "migrations"), nil
	default:
		return filepath.Join(base, service, "internal", "database", "migrations"), nil
	}
}
