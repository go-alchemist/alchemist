package setup

import "github.com/orochaa/go-clack/prompts"

type Feature struct {
	Name    string
	Enabled bool
}

func (m *model) AddFeatures() {
	featuresOrder := []string{
		"microservice_architecture",
	}

	featuresAvailable := map[string]string{
		"microservice_architecture": "Microservice Architecture",
	}

	options := make([]*prompts.MultiSelectOption[string], 0, len(featuresOrder))
	for _, key := range featuresOrder {
		label := featuresAvailable[key]
		options = append(options, &prompts.MultiSelectOption[string]{
			Value: key,
			Label: label,
		})
	}
	selectedFeatures, err := prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: "Select the features you want to enable:",
		Options: options,
	})
	prompts.ExitOnError(err)

	for key := range featuresAvailable {
		enabled := contains(selectedFeatures, key)
		m.Features = append(m.Features, Feature{
			Name:    key,
			Enabled: enabled,
		})
	}
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
