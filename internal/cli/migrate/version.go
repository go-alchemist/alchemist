package migrate

import "github.com/urfave/cli/v2"

func VersionMigration(c *cli.Context) error {
	return runMigrationCmd(c, "version")
}
