package migrate

import "github.com/urfave/cli/v2"

func ForceMigration(c *cli.Context) error {
	return runMigrationCmd(c, "force")
}
