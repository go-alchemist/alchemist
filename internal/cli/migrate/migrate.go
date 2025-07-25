package migrate

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-alchemist/alchemist/internal/cli/response"
)

var (
	cfgFile string
	RunCmd  = &cobra.Command{
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
	RunCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .env)")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(".env")
	}
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		response.Info(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}
}

func Migrate(cmd *cobra.Command, args []string) {
	_, dbURL, err := getDatabaseURL()
	if err != nil {
		response.Error("Error getting database URL: " + err.Error())
		return
	}

	dir, _ := cmd.Flags().GetString("dir")
	if dir == "" {
		dir = "internal/database/migrations"
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
	dbConn := viper.GetString("DB_CONNECTION")
	switch dbConn {
	case "sqlite":
		driver = "sqlite3"
		url = fmt.Sprintf("sqlite3://%s", viper.GetString("DB_DATABASE"))
	case "postgres":
		driver = "postgres"
		url = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			viper.GetString("DB_USERNAME"),
			viper.GetString("DB_PASSWORD"),
			viper.GetString("DB_HOST"),
			viper.GetString("DB_PORT"),
			viper.GetString("DB_DATABASE"),
		)
	case "mysql":
		driver = "mysql"
		url = fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
			viper.GetString("DB_USERNAME"),
			viper.GetString("DB_PASSWORD"),
			viper.GetString("DB_HOST"),
			viper.GetString("DB_PORT"),
			viper.GetString("DB_DATABASE"),
		)
	default:
		response.Error("Unsupported DB_CONNECTION or missing environment variables")
		return
	}

	return
}
