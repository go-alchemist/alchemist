package utils

import (
	"os"

	"github.com/orochaa/go-clack/prompts"

	"github.com/go-alchemist/alchemist/internal/cli/components"
)

func SelectModule(basePath string) string {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		prompts.Error(components.Red.Render("Could not read the modules directory. Check the path and permissions."))
		os.Exit(0)
		return ""
	}

	var modules []string
	for _, entry := range entries {
		if entry.IsDir() {
			modules = append(modules, entry.Name())
		}
	}

	if len(modules) == 0 {
		prompts.Error(components.Red.Render("No modules found in the directory."))
		os.Exit(0)
		return ""
	}

	var options []*prompts.SelectOption[string]
	for _, module := range modules {
		options = append(options, &prompts.SelectOption[string]{
			Label: module,
			Value: module,
		})
	}

	moduleDir, error := prompts.Select(prompts.SelectParams[string]{
		Message:  "Select a module:",
		Options:  options,
		Required: true,
	})
	if error != nil {
		prompts.ExitOnError(error)
		os.Exit(0)
		return ""
	}

	return moduleDir
}
