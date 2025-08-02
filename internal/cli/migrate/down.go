package migrate

import "github.com/urfave/cli/v2"

func RollbackMigration(c *cli.Context) error {
	return runMigrationCmd(c, "down")
}
