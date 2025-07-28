package components

import (
	"os"

	"github.com/orochaa/go-clack/prompts"
)

func SelectDomain(basePath string) string {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		prompts.Error(Red.Render("Could not read the domains directory. Check the path and permissions."))
		os.Exit(0)
		return ""

	}

	var domains []string
	for _, entry := range entries {
		if entry.IsDir() {
			domains = append(domains, entry.Name())
		}
	}

	if len(domains) == 0 {
		prompts.Error(Red.Render("No domains found in the directory."))
		os.Exit(0)
		return ""
	}

	var options []*prompts.SelectOption[string]
	for _, domain := range domains {
		options = append(options, &prompts.SelectOption[string]{
			Label: domain,
			Value: domain,
		})
	}

	domainDir, error := prompts.Select(prompts.SelectParams[string]{
		Message:  "Select a domain:",
		Options:  options,
		Required: true,
	})
	if error != nil {
		prompts.ExitOnError(error)
		os.Exit(0)
		return ""
	}

	return domainDir
}
