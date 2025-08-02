package setup

import "github.com/orochaa/go-clack/prompts"

func (m *model) SelectDefaultConfig() {
	configFiles := []string{
		".env",
		"config.yaml",
		"config.yml",
		"config.json",
		"config.toml",
		"config.ini",
	}

	options := make([]*prompts.SelectOption[string], 0, len(configFiles))
	for _, name := range configFiles {
		options = append(options, &prompts.SelectOption[string]{
			Value: name,
			Label: name,
		})
	}

	result, err := prompts.Select[string](prompts.SelectParams[string]{
		Message: "Select the main application configuration file:",
		Options: options,
	})

	prompts.ExitOnError(err)

	m.ConfigFile = result

}
