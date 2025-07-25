package response

import (
	"os"
	"path"
	"regexp"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"

	"github.com/go-alchemist/alchemist/internal/cli/config"
)

func Error(message string) {
	if message != "" {
		color.Red("🔴 %s", message)
		os.Exit(1)
	}
}

func Success(message string) {
	if message != "" {
		color.Green("🟢 %s", message)
	}
	os.Exit(0)
}

func Info(message string) {
	if message != "" {
		color.Blue("🔵 %s", message)
	}
}
func Warning(message string) {
	if message != "" {
		color.Yellow("🟡 %s", message)
	}
}

func SelectMicroserviceIfEnabled(baseDir string) string {
	enabled := config.Config.GetBool("features.microservice.enabled")
	if !enabled {
		Info("Microservice feature is not enabled.")
		return ""
	}

	serDir := config.Config.GetString("features.microservice.directory")
	if serDir == "" {
		serDir = "internal/services"
	}

	regexStr := config.Config.GetString("features.microservice.regex")
	if regexStr == "" {
		regexStr = ".*service.*"
	}
	re := regexp.MustCompile(regexStr)

	dir := path.Join(serDir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		Error("Directory not found: " + dir)
	}

	var services []string
	for _, entry := range entries {
		if entry.IsDir() && re.MatchString(entry.Name()) {
			services = append(services, entry.Name())
		}
	}

	if len(services) == 0 {
		Warning("No microservices found.")
		os.Exit(1)
		return ""
	}

	prompt := promptui.Select{
		Label:    "Select a microservice",
		Items:    services,
		HideHelp: false,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "🔹 {{ . | bold | cyan }}",
			Inactive: "   {{ . | faint }}",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		Warning("Selection cancelled")
		return ""
	}

	dirFinal := path.Join(serDir, result, baseDir)
	return dirFinal
}
