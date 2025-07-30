package make

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/orochaa/go-clack/prompts"
	"github.com/urfave/cli/v2"

	"github.com/go-alchemist/alchemist/internal/cli/templates"
	"github.com/go-alchemist/alchemist/internal/cli/utils"
)

func MakeRepository(c *cli.Context) error {
	name, err := getNamePrompt(c, "Repository")
	if err != nil {
		return err
	}
	base, structure, serviceDir, domainDir, modularDir := resolveDirs("repository", c)
	targetDir, err := RepositoryTargetDir(base, structure, serviceDir, domainDir, modularDir)
	if err != nil {
		return utils.PrintErrorAndReturn("Could not determine target directory for repository")
	}
	if err := utils.EnsureDir(targetDir); err != nil {
		return utils.PrintErrorAndReturn("Could not create the repository directory.")
	}
	filePath := filepath.Join(targetDir, strings.ToLower(name)+".go")
	if utils.FileExists(filePath) {
		return utils.PrintErrorAndReturn("A repository with this name already exists.")
	}
	tmpl, err := templates.GetRepositoryTemplate()
	if err != nil {
		return utils.PrintErrorAndReturn("Could not load the repository template.")
	}
	if err := writeSingleFile(filePath, tmpl, map[string]string{"RepositoryName": name}); err != nil {
		return utils.PrintErrorAndReturn("Could not generate the repository file.")
	}
	prompts.Success(fmt.Sprintf("Repository %s created successfully at:", name))
	utils.PrintSuccess(filePath)
	return nil
}

func RepositoryTargetDir(base, structure, service, domain, module string) (string, error) {
	switch structure {
	case "clean_architecture":
		return filepath.Join(base, service, "domain", "repository"), nil
	case "domain_driven":
		return filepath.Join(base, service, domain, "repository"), nil
	case "modular":
		return filepath.Join(base, service, "modules", module, "repository"), nil
	case "layered":
		return filepath.Join(base, service, "repository"), nil
	case "custom":
		return utils.SelectCustomDirectory(filepath.Join(base, service), "repository"), nil
	default:
		return filepath.Join(base, service, "internal", "repository"), nil
	}
}
