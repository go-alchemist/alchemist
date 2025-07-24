package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/go-alchemist/alchemist/cmd/make"
	"github.com/go-alchemist/alchemist/cmd/migrate"
)

var version = "v1.0.0"

var rootCmd = &cobra.Command{
	Use:     "alchemist",
	Short:   "Alchemist - CLI to projects in Go",
	Long:    `Alchemist is a CLI tool to manage and interact with Go projects.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(make.MakeCmd)
	rootCmd.AddCommand(migrate.MigrateCmd)
	rootCmd.SetVersionTemplate("Alchemist CLI version: {{.Version}}\n")

}
