package components

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/orochaa/go-clack/prompts"

	"github.com/go-alchemist/alchemist/internal/cli/config"
)

func SelectMicroserviceIfEnabled(structure string) string {
	enabled := config.Config.GetBool("features.microservice_architecture.enabled")
	if !enabled {
		return ""
	}

	serDir := config.Config.GetString("features.microservice_architecture.directory")
	if serDir == "" {
		serDir = "./modules"
	}

	regexStr := config.Config.GetString("features.microservice_architecture.regex")
	if regexStr == "" {
		regexStr = ".*service.*"
	}
	re := regexp.MustCompile(regexStr)

	dir := path.Join(serDir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		prompts.Error(Red.Render("Error performing operation"))
		fmt.Println((Red.Render("Cannot read the microservices directory. Check the path and permissions.")))
		os.Exit(0)
		return ""
	}

	var services []string
	for _, entry := range entries {
		if entry.IsDir() && re.MatchString(entry.Name()) {
			services = append(services, entry.Name())
		}
	}

	if len(services) == 0 {
		prompts.Error(Red.Render("Error performing operation"))
		fmt.Println(Red.Render("No microservices found matching the regex: " + regexStr))
		os.Exit(0)
		return ""
	}

	var options []*prompts.SelectOption[string]
	for _, svc := range services {
		options = append(options, &prompts.SelectOption[string]{
			Label: svc,
			Value: svc,
		})
	}

	result, err := prompts.Select(prompts.SelectParams[string]{
		Message: "Select a microservice",
		Options: options,
	})
	if err != nil {
		if prompts.IsCancel(err) {
			prompts.Error(Red.Render("Selection cancelled"))
			os.Exit(0)
			return ""
		}
	}

	dir = path.Join(serDir, result)

	return dir
}
