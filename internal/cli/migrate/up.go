package migrate

import "github.com/urfave/cli/v2"

func RunMigration(c *cli.Context) error {
	return runMigrationCmd(c, "up")
}
