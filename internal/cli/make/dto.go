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

func MakeDTO(c *cli.Context) error {
	name, err := getNamePrompt(c, "DTO")
	if err != nil {
		return err
	}
	base, structure, serviceDir, domainDir, modularDir := resolveDirs("dto", c)
	targetDir, err := DTOTargetDir(base, structure, serviceDir, domainDir, modularDir)
	if err != nil {
		return utils.PrintErrorAndReturn("Could not determine target directory for DTO")
	}
	if err := utils.EnsureDir(targetDir); err != nil {
		return utils.PrintErrorAndReturn("Could not create the DTO directory.")
	}
	filePath := filepath.Join(targetDir, strings.ToLower(name)+".go")
	if utils.FileExists(filePath) {
		return utils.PrintErrorAndReturn("A DTO with this name already exists.")
	}
	tmpl, err := templates.GetDTOTemplate()
	if err != nil {
		return utils.PrintErrorAndReturn("Could not load the DTO template.")
	}
	if err := writeSingleFile(filePath, tmpl, map[string]string{"DTOName": name}); err != nil {
		return utils.PrintErrorAndReturn("Could not generate the DTO file.")
	}
	prompts.Success(fmt.Sprintf("DTO %s created successfully at:", name))
	utils.PrintSuccess(filePath)
	return nil
}

func DTOTargetDir(base, structure, service, domain, module string) (string, error) {
	switch structure {
	case "clean_architecture":
		return filepath.Join(base, service, "application", "dto"), nil
	case "domain_driven":
		return filepath.Join(base, service, domain, "dto"), nil
	case "modular":
		return filepath.Join(base, service, "modules", module, "dto"), nil
	case "layered":
		return filepath.Join(base, service, "dto"), nil
	case "custom":
		return utils.SelectCustomDirectory(filepath.Join(base, service), "dto"), nil
	default:
		return filepath.Join(base, service, "internal", "dto"), nil
	}
}
