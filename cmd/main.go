package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/go-alchemist/alchemist/internal/cli/setup"
)

var Version = "v1.0.0"

func main() {
	app := &cli.App{
		Name:    "alchemist",
		Usage:   "Alchemist - CLI to projects in Go",
		Version: Version,
		Commands: []*cli.Command{
			{
				Name:   "setup",
				Usage:  "Run the Alchemist setup",
				Action: setup.RunSetup,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
