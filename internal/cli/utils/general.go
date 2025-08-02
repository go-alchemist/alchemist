package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/orochaa/go-clack/prompts"

	"github.com/go-alchemist/alchemist/internal/cli/components"
)

var configFiles = []string{
	".env",
	"config.yaml",
	"config.yml",
	"config.json",
	"config.toml",
	"config.ini",
}

func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func PrintErrorAndReturn(msg string) error {
	prompts.Error(components.Red.Render("Error performing operation"))
	prompts.Outro(components.Red.Render(msg))
	return nil
}

func PrintSuccess(msg string) {
	prompts.Outro(components.Green.Render(msg))
}

func PrintSuccessf(format string, args ...interface{}) {
	fmt.Printf(components.Green.Render(format), args...)
}

func PrintErrorAndReturnF(format string, args ...interface{}) {
	prompts.Error(components.Red.Render("Error performing operation"))
	fmt.Printf(components.Red.Render(format), args...)
	os.Exit(0)
}

func FindConfigFile(serviceDir string) (string, string, error) {
	for _, file := range configFiles {
		path := filepath.Join(serviceDir, file)
		if _, err := os.Stat(path); err == nil {
			ext := strings.TrimPrefix(filepath.Ext(file), ".")
			if ext == "" {
				ext = "env"
			}
			return path, ext, nil
		}
	}
	return "", "", fmt.Errorf("No config file found in %s", serviceDir)
}
