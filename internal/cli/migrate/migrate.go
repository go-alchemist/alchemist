package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/orochaa/go-clack/prompts"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"github.com/go-alchemist/alchemist/internal/cli/config"
	makeCommand "github.com/go-alchemist/alchemist/internal/cli/make"
	"github.com/go-alchemist/alchemist/internal/cli/utils"
)

func runMigrationCmd(c *cli.Context, mode string) error {
	base := "."
	structure := utils.DetectProjectStructure(base)
	serviceDir := utils.SelectMicroserviceIfEnabled(structure)
	if serviceDir == "" {
		utils.PrintErrorAndReturnF("Microservice not found or not enabled")
	}

	migrationsDir := c.String("dir")
	targetDir, err := utils.GetTargetDir(base, structure, serviceDir, migrationsDir, makeCommand.MigrationTargetDir)
	if err != nil {
		utils.PrintErrorAndReturnF("Could not determine target directory for migration")
	}
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		utils.PrintErrorAndReturnF(fmt.Sprintf("Migration directory does not exist: %s", targetDir))
	}

	cfgPath, cfgType, err := utils.FindConfigFile(serviceDir)
	if err != nil {
		return utils.PrintErrorAndReturn(fmt.Sprintf("No config file found in microservice: %s", serviceDir))
	}

	viper.SetConfigFile(cfgPath)
	viper.SetConfigType(cfgType)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return utils.PrintErrorAndReturn(fmt.Sprintf("Error reading config: %s", err.Error()))
	}
	_, dbURL, err := getDatabaseURL()
	if err != nil {
		utils.PrintErrorAndReturnF("Error getting database URL: " + err.Error())
	}

	dbConfig, err := getConfigDB()
	if err != nil {
		utils.PrintErrorAndReturnF("Error getting database configuration: " + err.Error())
	}
	schema := ""
	if dbConfig["driver"] == "postgres" {
		schema = dbConfig["schema"]
		if err := EnsurePostgresSchemaExists(dbURL, schema); err != nil {
			return utils.PrintErrorAndReturn("Error ensuring schema exists: " + err.Error())
		}
	}

	m, err := migrate.New("file://"+targetDir, dbURL)
	if err != nil {
		utils.PrintErrorAndReturnF("Error initializing migration: " + err.Error())
	}
	defer m.Close()

	switch mode {
	case "up":
		if err := m.Up(); err != nil && err.Error() != "no change" {
			utils.PrintErrorAndReturnF("Error applying migrations: " + err.Error())
		}
		prompts.Success("All migrations applied successfully!")
	case "down":
		steps := c.Int("steps")
		if steps == 0 {
			steps = 1
		}
		err := m.Steps(-steps)
		if err != nil {
			if err.Error() == "no change" {
				prompts.Info("No migrations to rollback.")
			} else if strings.Contains(err.Error(), "file does not exist") {
				prompts.Info("No more migrations to rollback.")
			} else {
				utils.PrintErrorAndReturnF(fmt.Sprintf("Error rolling back %d migration(s): %v", steps, err))
			}
		} else {
			prompts.Success(fmt.Sprintf("Rolled back %d migration(s) successfully!", steps))
		}
	case "force":
		versionStr := c.String("version")
		if versionStr == "" {
			utils.PrintErrorAndReturnF("You must provide a --version for force")
		}
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			utils.PrintErrorAndReturnF("Invalid version: " + versionStr)
		}
		if err := m.Force(version); err != nil {
			utils.PrintErrorAndReturnF(fmt.Sprintf("Error forcing migration version: %v", err))
		}
		prompts.Success(fmt.Sprintf("Forced migration version to %d", version))
	case "version":
		v, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			utils.PrintErrorAndReturnF(fmt.Sprintf("Error getting migration version: %v", err))
		}
		status := "clean"
		if dirty {
			status = "dirty"
		}
		msg := fmt.Sprintf("Current migration version: %d (%s)", v, status)
		if err == migrate.ErrNilVersion {
			msg = "No migration has been applied yet."
		}
		prompts.Info(msg)
	default:
		utils.PrintErrorAndReturnF("Unknown migration command")
	}

	return nil
}

func getDatabaseURL() (driver string, url string, err error) {
	dbConfig := config.Config.GetStringMapString("config.database")
	if len(dbConfig) == 0 {
		err = errors.New("Missing database configuration in config file")
		return
	}

	get := func(key string) string {
		if envKey, ok := dbConfig[key]; ok {
			return viper.GetString(envKey)
		}
		return ""
	}

	dbDriver := get("driver")
	switch dbDriver {
	case "sqlite", "sqlite3":
		driver = "sqlite3"
		url = fmt.Sprintf("sqlite3://%s", get("name"))
	case "postgres":
		driver = "postgres"
		schema := get("schema")
		searchPath := ""
		if schema != "" {
			searchPath = fmt.Sprintf("&search_path=%s", schema)
		}
		url = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable%s",
			get("user"), get("password"), get("host"), get("port"), get("name"),
			searchPath,
		)
	case "mysql":
		driver = "mysql"
		url = fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
			get("user"), get("password"), get("host"), get("port"), get("name"),
		)
	default:
		err = errors.New("unsupported DB driver or missing environment variables")
	}
	return
}

func getConfigDB() (map[string]string, error) {
	dbConfig := config.Config.GetStringMapString("config.database")
	if len(dbConfig) == 0 {
		return nil, errors.New("Missing database configuration in config file")
	}

	result := make(map[string]string)
	for key, envKey := range dbConfig {
		result[key] = viper.GetString(envKey)
	}
	return result, nil
}

func EnsurePostgresSchemaExists(dbURL, schema string) error {
	if schema == "" {
		return nil
	}

	connStr := dbURL
	if idx := strings.Index(connStr, "search_path="); idx != -1 {
		connStr = connStr[:idx-1]
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1);`
	err = db.QueryRow(query, schema).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if schema exists: %w", err)
	}

	if !exists {
		_, err := db.Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schema))
		if err != nil {
			return fmt.Errorf("failed to create schema '%s': %w", schema, err)
		}
	}
	return nil
}
