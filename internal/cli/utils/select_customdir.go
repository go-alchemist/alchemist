package utils

import (
	"path/filepath"

	"github.com/go-alchemist/alchemist/internal/cli/config"
)

func SelectCustomDirectory(basePath, componentType string) string {
	customDir := config.Config.GetString("custom_path." + componentType)
	return filepath.Join(basePath, customDir)
}
