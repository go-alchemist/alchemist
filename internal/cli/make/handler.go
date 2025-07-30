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

func MakeHandler(c *cli.Context) error {
	handlerName, err := getNamePrompt(c, "Handler")
	if err != nil {
		return err
	}

	base := "."
	structure := utils.DetectProjectStructure(base)

	serviceDir := utils.SelectMicroserviceIfEnabled(structure)
	if serviceDir == "" {
		return utils.PrintErrorAndReturn("Microservice not found or not enabled")
	}

	var domainDir, modularDir string
	switch structure {
	case "modular":
		modularDir = utils.SelectModule(filepath.Join(base, serviceDir, "modules"))
	case "domain_driven":
		domainDir = utils.SelectDomain(filepath.Join(base, serviceDir))
	}

	targetDir, err := HandlerTargetDir(base, structure, serviceDir, domainDir, modularDir)
	if err != nil {
		return utils.PrintErrorAndReturn("Could not determine target directory for handler")
	}

	if err := utils.EnsureDir(targetDir); err != nil {
		return utils.PrintErrorAndReturn("Could not create the handler directory. Please check your permissions and try again.")
	}

	filePath := filepath.Join(targetDir, strings.ToLower(handlerName)+".go")
	if utils.FileExists(filePath) {
		return utils.PrintErrorAndReturn("A handler with this name already exists. Choose another name.")
	}

	tmpl, err := templates.GetHandlerTemplate()
	if err != nil {
		return utils.PrintErrorAndReturn("Could not load the handler template.")
	}

	if err := writeSingleFile(filePath, tmpl, map[string]string{"HandlerName": handlerName}); err != nil {
		return utils.PrintErrorAndReturn("Could not generate the handler file from the template. Please check the template syntax.")
	}

	prompts.Success(fmt.Sprintf("Handler %s created successfully at:", handlerName))

	utils.PrintSuccess(filePath)

	return nil
}

func HandlerTargetDir(base, structure, service, domain, module string) (string, error) {
	switch structure {
	case "clean_architecture":
		return filepath.Join(base, service, "application", "service"), nil
	case "domain_driven":
		return filepath.Join(base, service, domain, "handler"), nil
	case "modular":
		return filepath.Join(base, service, "modules", module, "handler"), nil
	case "layered":
		return filepath.Join(base, service, "handler"), nil
	case "custom":
		return utils.SelectCustomDirectory(filepath.Join(base, service), "handler"), nil
	default:
		return filepath.Join(base, service, "internal", "handler"), nil
	}
}
