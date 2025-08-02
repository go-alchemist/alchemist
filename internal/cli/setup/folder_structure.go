package setup

import "github.com/orochaa/go-clack/prompts"

func (m *model) SelectFolderStructure() {
	structureOrder := []string{
		"standard_layout",
		"domain_driven",
		"layered",
		"clean_architecture",
		"modular",
		"custom",
	}
	structureAvailable := map[string]string{
		"standard_layout":    "Standard/Idiomatic (cmd/pkg/internal)",
		"domain_driven":      "Domain-driven Design (by business feature/domain)",
		"layered":            "Layered (handler/service/repository)",
		"clean_architecture": "Clean Architecture/Hexagonal",
		"modular":            "Modular/Monorepo (multi-module, go.work)",
		"custom":             "Custom...",
	}

	options := make([]*prompts.SelectOption[string], 0, len(structureOrder))
	for _, key := range structureOrder {
		label := structureAvailable[key]
		options = append(options, &prompts.SelectOption[string]{
			Value: key,
			Label: label,
		})
	}

	folderStructure, err := prompts.Select(prompts.SelectParams[string]{
		Message: "Pick a project structure:",
		Options: options,
	})
	prompts.ExitOnError(err)

	m.FolderStructure = folderStructure

}
