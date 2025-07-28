package components

import (
	"github.com/go-alchemist/alchemist/internal/cli/config"
)

func DetectProjectStructure(base string) string {
	switch config.Config.GetString("path_structure") {
	case "clean_architecture":
		return "clean_architecture"
	case "modular":
		return "modular"
	case "domain_driven":
		return "domain_driven"
	case "layered":
		return "layered"
	case "standard_layout":
		return "standard_layout"
	case "custom":
		return "custom"
	}

	return "standard_layout"
}
