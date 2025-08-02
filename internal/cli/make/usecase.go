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

func MakeUsecase(c *cli.Context) error {
	name, err := getNamePrompt(c, "Usecase")
	if err != nil {
		return err
	}
	base, structure, serviceDir, domainDir, modularDir := resolveDirs("usecase", c)
	sDir := c.String("dir")
	targetDir, err := utils.GetTargetDirWithDomainModule(
		serviceDir, domainDir, modularDir, sDir, UsecaseTargetDir, base, structure,
	)
	if err != nil {
		return utils.PrintErrorAndReturn("Could not determine target directory for DTO")
	}
	if err := utils.EnsureDir(targetDir); err != nil {
		return utils.PrintErrorAndReturn("Could not create the usecase directory.")
	}
	filePath := filepath.Join(targetDir, strings.ToLower(name)+".go")
	if utils.FileExists(filePath) {
		return utils.PrintErrorAndReturn("A usecase with this name already exists.")
	}
	tmpl, err := templates.GetUsecaseTemplate()
	if err != nil {
		return utils.PrintErrorAndReturn("Could not load the usecase template.")
	}
	if err := writeSingleFile(filePath, tmpl, map[string]string{"UsecaseName": name}); err != nil {
		return utils.PrintErrorAndReturn("Could not generate the usecase file.")
	}
	prompts.Success(fmt.Sprintf("Usecase %s created successfully at:", name))
	utils.PrintSuccess(filePath)
	return nil
}

func UsecaseTargetDir(base, structure, service, domain, module string) (string, error) {
	switch structure {
	case "clean_architecture":
		return filepath.Join(base, service, "application", "usecase"), nil
	case "domain_driven":
		return filepath.Join(base, service, domain, "usecase"), nil
	case "modular":
		return filepath.Join(base, service, "modules", module, "usecase"), nil
	case "layered":
		return filepath.Join(base, service, "usecase"), nil
	case "custom":
		return utils.SelectCustomDirectory(filepath.Join(base, service), "usecase"), nil
	default:
		return filepath.Join(base, service, "internal", "usecase"), nil
	}
}
