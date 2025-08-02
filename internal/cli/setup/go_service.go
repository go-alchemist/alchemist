package setup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var goVersion = "1.24.5"

func createGoMod(basePath, moduleName string) error {
	goModPath := filepath.Join(basePath, "go.mod")
	content := fmt.Sprintf("module %s\n\ngo %s\n", moduleName, goVersion)
	return os.WriteFile(goModPath, []byte(content), 0644)
}

func createGoWork(root string, modulePaths []string) error {
	goWorkPath := filepath.Join(root, "go.work")
	content := fmt.Sprintf("go %s\n\nuse (\n", goVersion)
	for _, mod := range modulePaths {
		content += fmt.Sprintf("\t%s\n", mod)
	}
	content += ")\n"
	return os.WriteFile(goWorkPath, []byte(content), 0644)
}

func createGoFilesWithFeatures(projectName string, features map[string]bool) error {
	root := projectName

	if features["microservice_architecture"] {
		modulesDir := filepath.Join(root, "modules")
		entries, err := ioutil.ReadDir(modulesDir)
		if err != nil {
			return fmt.Errorf("erro ao ler a pasta de módulos: %w", err)
		}

		modulePaths := []string{}
		for _, entry := range entries {
			if entry.IsDir() {
				modName := entry.Name()
				moduleGoName := GoModuleNameFromFolder(modName)
				modPath := filepath.Join(modulesDir, modName)

				if err := createGoMod(modPath, moduleGoName); err != nil {
					return err
				}
				modulePaths = append(modulePaths, fmt.Sprintf("./modules/%s", modName))
			}
		}

		return createGoWork(root, modulePaths)
	}

	return createGoMod(root, projectName)
}

func GoModuleNameFromFolder(folder string) string {
	var result strings.Builder
	capitalize := false
	for i, r := range folder {
		if r == '_' || r == '-' {
			capitalize = true
			continue
		}
		if capitalize {
			result.WriteRune(unicode.ToUpper(r))
			capitalize = false
		} else if i == 0 {
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
