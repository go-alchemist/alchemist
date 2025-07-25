package migrate

import (
	"fmt"
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-alchemist/alchemist/internal/cli/config"
	"github.com/go-alchemist/alchemist/internal/cli/response"
)

var (
	RunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run all pending SQL migrations",
		Run:   Migrate,
	}
)

func Execute() error {
	return RunCmd.Execute()
}

func init() {
	RunCmd.Flags().String("dir", "internal/database/migrations", "Directory for migrations")
}

func Migrate(cmd *cobra.Command, args []string) {
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

	envPath := path.Join(dir, config.Config.GetString("env_config.file"))
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		response.Error(fmt.Sprintf("No .env found in microservice: %s", envPath))
		return
	}

	viper.SetConfigFile(envPath)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		response.Error("Error reading .env: " + err.Error())
		return
	}
	response.Info(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))

	_, dbURL, err := getDatabaseURL()
	if err != nil {
		response.Error("Error getting database URL: " + err.Error())
		return
	}

	m, err := migrate.New(
		"file://"+dir,
		dbURL,
	)
	if err != nil {
		response.Error("Error initializing migration: " + err.Error())
		return
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		response.Error("Error applying migrations: " + err.Error())
		return
	}

	response.Success("All migrations applied successfully!")
}

func getDatabaseURL() (driver string, url string, err error) {
	dbConfig := config.Config.GetStringMapString("env_config.database")

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
		url = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			get("user"),
			get("password"),
			get("host"),
			get("port"),
			get("name"),
		)
	case "mysql":
		driver = "mysql"
		url = fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
			get("user"),
			get("password"),
			get("host"),
			get("port"),
			get("name"),
		)
	default:
		response.Error("Unsupported DB driver or missing environment variables")
		return
	}

	return
}
