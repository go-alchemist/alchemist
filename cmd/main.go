package main

import (
	"log"
	"os"

	"github.com/orochaa/go-clack/prompts"
	"github.com/urfave/cli/v2"

	"github.com/go-alchemist/alchemist/internal/cli/components"
	"github.com/go-alchemist/alchemist/internal/cli/config"
	"github.com/go-alchemist/alchemist/internal/cli/make"
	"github.com/go-alchemist/alchemist/internal/cli/migrate"
	"github.com/go-alchemist/alchemist/internal/cli/setup"
)

var Version = "v1.0.0"

func main() {
	app := &cli.App{
		Name:     "alchemist",
		HelpName: "Alchemist",
		Usage:    "CLI to projects in Gos",
		Version:  Version,
		Before: func(c *cli.Context) error {
			prompts.Intro(components.Banner())
			return nil
		},
		CommandNotFound: func(c *cli.Context, command string) {
			prompts.Error(components.Red.Render("Command not found"))
			os.Exit(1)
		},
		Commands: []*cli.Command{
			{
				Name:   "setup",
				Usage:  "Run the Alchemist setup",
				Action: setup.RunSetup,
			},
			{
				Name:  "make",
				Usage: "Generate boilerplate code for various components",
				Subcommands: []*cli.Command{
					{
						Name:   "handler",
						Usage:  "Create a new handler",
						Action: make.MakeHandler,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "dir",
								Aliases: []string{"d"},
								Usage:   "Directory to create the handler in",
							},
						},
						Before: initConfig,
					},
					{
						Name:  "model",
						Usage: "Create a new model",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "dir",
								Aliases: []string{"d"},
								Usage:   "Directory to create the model in",
							},
						},
						Action: make.MakeModel,
						Before: initConfig,
					},
					{
						Name:  "repository",
						Usage: "Create a new repository",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "dir",
								Aliases: []string{"d"},
								Usage:   "Directory to create the repository in",
							},
						},
						Action: make.MakeRepository,
						Before: initConfig,
					},
					{
						Name:  "service",
						Usage: "Create a new service",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "dir",
								Aliases: []string{"d"},
								Usage:   "Directory to create the service in",
							},
						},
						Action: make.MakeService,
						Before: initConfig,
					},
					{
						Name:  "migration",
						Usage: "Create a new migration",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "dir",
								Aliases: []string{"d"},
								Usage:   "Directory to create the migration in",
							},
						},
						Action: make.MakeMigration,
						Before: initConfig,
					},

					{
						Name:  "dto",
						Usage: "Create a new DTO",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "dir",
								Aliases: []string{"d"},
								Usage:   "Directory to create the DTO in",
							},
						},
						Action: make.MakeDTO,
						Before: initConfig,
					},
					{
						Name:  "usecase",
						Usage: "Create a new usecase",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "dir",
								Aliases: []string{"d"},
								Usage:   "Directory to create the usecase in",
							},
						},
						Action: make.MakeUsecase,
						Before: initConfig,
					},

					{
						Name:  "test",
						Usage: "Create a new test file",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "dir",
								Aliases: []string{"d"},
								Usage:   "Directory to create the test in",
							},
						},
						Action: make.MakeTest,
						Before: initConfig,
					},
				},
			},
			{
				Name:  "migrate",
				Usage: "Run migration commands for the selected microservice",
				Subcommands: []*cli.Command{
					{
						Name:  "up",
						Usage: "Apply all up migrations",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "dir", Aliases: []string{"d"}, Usage: "Directory for migrations"},
						},
						Action: migrate.RunMigration,
						Before: initConfig,
					},
					{
						Name:  "down",
						Usage: "Rollback the last or N migrations",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "dir", Aliases: []string{"d"}, Usage: "Directory for migrations"},
							&cli.IntFlag{Name: "steps", Aliases: []string{"s"}, Usage: "Number of migrations to rollback"},
						},
						Action: migrate.RollbackMigration,
						Before: initConfig,
					},
					{
						Name:  "force",
						Usage: "Force set migration version",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "dir", Aliases: []string{"d"}, Usage: "Directory for migrations"},
							&cli.StringFlag{Name: "version", Aliases: []string{"v"}, Usage: "Version to force set"},
						},
						Action: migrate.ForceMigration,
						Before: initConfig,
					},
					{
						Name:  "version",
						Usage: "Show current migration version",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "dir", Aliases: []string{"d"}, Usage: "Directory for migrations"},
						},
						Action: migrate.VersionMigration,
						Before: initConfig,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initConfig(c *cli.Context) error {
	config.InitConfig()
	return nil
}
