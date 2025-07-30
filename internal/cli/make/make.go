package make

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/orochaa/go-clack/prompts"
	"github.com/urfave/cli/v2"

	"github.com/go-alchemist/alchemist/internal/cli/utils"
)

func getNamePrompt(c *cli.Context, kind string) (string, error) {
	if c.NArg() > 0 && c.Args().Get(0) != "" {
		return c.Args().Get(0), nil
	}
	return prompts.Text(prompts.TextParams{
		Message:  kind + " name:",
		Required: true,
		Validate: func(value string) error {
			switch {
			case value == "":
				return fmt.Errorf("%s name cannot be empty", kind)
			case strings.Contains(value, " "):
				return fmt.Errorf("%s name cannot contain spaces", kind)
			}
			return nil
		},
	})
}

func writeSingleFile(filePath string, tmpl *template.Template, data map[string]string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}

func resolveDirs(kind string, c *cli.Context) (base, structure, serviceDir, domainDir, modularDir string) {
	base = "."
	structure = utils.DetectProjectStructure(base)
	serviceDir = utils.SelectMicroserviceIfEnabled(structure)
	switch structure {
	case "modular":
		modularDir = utils.SelectModule(filepath.Join(base, serviceDir, "modules"))
	case "domain_driven":
		domainDir = utils.SelectDomain(filepath.Join(base, serviceDir))
	}
	return
}
