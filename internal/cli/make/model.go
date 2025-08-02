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

func MakeModel(c *cli.Context) error {
	modelName, err := getNamePrompt(c, "Model")
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

	sDir := c.String("dir")
	targetDir, err := utils.GetTargetDirWithDomainModule(
		serviceDir, domainDir, modularDir, sDir, ModelTargetDir, base, structure,
	)
	if err != nil {
		return utils.PrintErrorAndReturn("Could not determine target directory for DTO")
	}

	if err := utils.EnsureDir(targetDir); err != nil {
		return utils.PrintErrorAndReturn("Could not create the model directory. Please check your permissions and try again.")
	}

	filePath := filepath.Join(targetDir, strings.ToLower(modelName)+".go")
	if utils.FileExists(filePath) {
		return utils.PrintErrorAndReturn("A model with this name already exists. Choose another name.")
	}

	tmpl, err := templates.GetModelTemplate()
	if err != nil {
		return utils.PrintErrorAndReturn("Could not load the model template.")
	}

	if err := writeSingleFile(filePath, tmpl, map[string]string{"ModelName": modelName}); err != nil {
		return utils.PrintErrorAndReturn("Could not generate the model file from the template. Please check the template syntax.")
	}

	prompts.Success(fmt.Sprintf("Model %s created successfully at:", modelName))
	utils.PrintSuccess(filePath)
	return nil
}

func ModelTargetDir(base, structure, service, domain, module string) (string, error) {
	switch structure {
	case "clean_architecture":
		return filepath.Join(base, service, "domain", "model"), nil
	case "domain_driven":
		return filepath.Join(base, service, domain, "model"), nil
	case "modular":
		return filepath.Join(base, service, "modules", module, "model"), nil
	case "layered":
		return filepath.Join(base, service, "model"), nil
	case "custom":
		return utils.SelectCustomDirectory(filepath.Join(base, service), "model"), nil
	default:
		return filepath.Join(base, service, "internal", "model"), nil
	}
}
