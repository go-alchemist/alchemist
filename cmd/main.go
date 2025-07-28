package main

import (
	"log"
	"os"

	"github.com/orochaa/go-clack/prompts"
	"github.com/urfave/cli/v2"

	"github.com/go-alchemist/alchemist/internal/cli/components"
	"github.com/go-alchemist/alchemist/internal/cli/config"
	"github.com/go-alchemist/alchemist/internal/cli/make"
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
				Name: "make",
				Subcommands: []*cli.Command{
					{
						Name:   "handler",
						Usage:  "Create a new handler",
						Action: make.MakeHandler,
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
