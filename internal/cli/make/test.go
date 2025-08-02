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

func MakeTest(c *cli.Context) error {
	name, err := getNamePrompt(c, "Test")
	if err != nil {
		return err
	}
	base, structure, serviceDir, domainDir, modularDir := resolveDirs("test", c)
	sDir := c.String("dir")
	targetDir, err := utils.GetTargetDirWithDomainModule(
		serviceDir, domainDir, modularDir, sDir, TestTargetDir, base, structure,
	)
	if err != nil {
		return utils.PrintErrorAndReturn("Could not determine target directory for DTO")
	}
	if err := utils.EnsureDir(targetDir); err != nil {
		return utils.PrintErrorAndReturn("Could not create the test directory.")
	}
	filePath := filepath.Join(targetDir, strings.ToLower(name)+"_test.go")
	if utils.FileExists(filePath) {
		return utils.PrintErrorAndReturn("A test file with this name already exists.")
	}
	tmpl, err := templates.GetTestTemplate()
	if err != nil {
		return utils.PrintErrorAndReturn("Could not load the test template.")
	}
	if err := writeSingleFile(filePath, tmpl, map[string]string{"TestName": name}); err != nil {
		return utils.PrintErrorAndReturn("Could not generate the test file.")
	}
	prompts.Success(fmt.Sprintf("Test %s created successfully at:", name))
	utils.PrintSuccess(filePath)
	return nil
}

func TestTargetDir(base, structure, service, domain, module string) (string, error) {
	switch structure {
	case "clean_architecture":
		return filepath.Join(base, service, "tests"), nil
	case "domain_driven":
		return filepath.Join(base, service, domain, "tests"), nil
	case "modular":
		return filepath.Join(base, service, "modules", module, "tests"), nil
	case "layered":
		return filepath.Join(base, service, "tests"), nil
	case "custom":
		return utils.SelectCustomDirectory(filepath.Join(base, service), "tests"), nil
	default:
		return filepath.Join(base, service, "internal", "tests"), nil
	}
}
