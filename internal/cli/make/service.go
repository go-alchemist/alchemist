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

func MakeService(c *cli.Context) error {
	name, err := getNamePrompt(c, "Service")
	if err != nil {
		return err
	}
	base, structure, serviceDir, domainDir, modularDir := resolveDirs("service", c)
	sDir := c.String("dir")
	targetDir, err := utils.GetTargetDirWithDomainModule(
		serviceDir, domainDir, modularDir, sDir, ServiceTargetDir, base, structure,
	)
	if err != nil {
		return utils.PrintErrorAndReturn("Could not determine target directory for DTO")
	}
	if err := utils.EnsureDir(targetDir); err != nil {
		return utils.PrintErrorAndReturn("Could not create the service directory.")
	}
	filePath := filepath.Join(targetDir, strings.ToLower(name)+".go")
	if utils.FileExists(filePath) {
		return utils.PrintErrorAndReturn("A service with this name already exists.")
	}
	tmpl, err := templates.GetServiceTemplate()
	if err != nil {
		return utils.PrintErrorAndReturn("Could not load the service template.")
	}
	if err := writeSingleFile(filePath, tmpl, map[string]string{"ServiceName": name}); err != nil {
		return utils.PrintErrorAndReturn("Could not generate the service file.")
	}
	prompts.Success(fmt.Sprintf("Service %s created successfully at:", name))
	utils.PrintSuccess(filePath)
	return nil
}

func ServiceTargetDir(base, structure, service, domain, module string) (string, error) {
	switch structure {
	case "clean_architecture":
		return filepath.Join(base, service, "application", "service"), nil
	case "domain_driven":
		return filepath.Join(base, service, domain, "service"), nil
	case "modular":
		return filepath.Join(base, service, "modules", module, "service"), nil
	case "layered":
		return filepath.Join(base, service, "service"), nil
	case "custom":
		return utils.SelectCustomDirectory(filepath.Join(base, service), "service"), nil
	default:
		return filepath.Join(base, service, "internal", "service"), nil
	}
}
