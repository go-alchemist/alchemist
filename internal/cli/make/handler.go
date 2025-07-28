package make

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/orochaa/go-clack/prompts"
	"github.com/urfave/cli/v2"

	"github.com/go-alchemist/alchemist/internal/cli/components"
	"github.com/go-alchemist/alchemist/internal/cli/templates"
)

func MakeHandler(c *cli.Context) error {
	var handlerName string

	if c.NArg() > 0 && c.Args().Get(0) != "" {
		handlerName = c.Args().Get(0)
	} else {
		name, err := prompts.Text(prompts.TextParams{
			Message:  "Handler name:",
			Required: true,
			Validate: func(value string) error {
				if value == "" {
					return errors.New("Handler name cannot be empty")
				}
				if strings.Contains(value, " ") {
					return errors.New("Handler name cannot contain spaces")
				}

				return nil
			},
		})
		if err != nil {
			prompts.ExitOnError(err)
			return err
		}
		handlerName = name
	}

	base := "."
	structure := components.DetectProjectStructure(base)

	serDir := components.SelectMicroserviceIfEnabled(structure)
	domainDir := ""
	modularDir := ""

	switch structure {
	case "modular":
		serviceDir := filepath.Join(base, serDir, "modules")
		modularDir = components.SelectModule(serviceDir)
	case "domain_driven":
		serviceDir := filepath.Join(base, serDir)
		domainDir = components.SelectDomain(serviceDir)
	}

	dir, err := HandlerTargetDir(base, structure, serDir, domainDir, modularDir, handlerName)
	if err != nil {
		prompts.Error(components.Red.Render("Error performing operation"))
		fmt.Println(components.Red.Render("\n  Could not determine target directory for handler"))
		return err
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			prompts.Error(components.Red.Render("Error performing operation"))
			fmt.Println(components.Red.Render("\n  Could not create the handler directory. Please check your permissions and try again."))
			os.Exit(0)
			return err
		}
	}

	filePath := fmt.Sprintf("%s/%s.go", dir, strings.ToLower(handlerName))
	if _, err := os.Stat(filePath); err == nil {
		prompts.Error(components.Red.Render("Error performing operation"))
		fmt.Println(components.Red.Render("\n  A handler with this name already exists. Choose another name."))
		return nil
	}
	tmpl, err := templates.GetHandlerTemplate()
	if err != nil {
		prompts.Error(components.Red.Render("Error performing operation"))
		fmt.Println(components.Red.Render("\n  Could not load the handler template."))
		os.Exit(0)
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		prompts.Error("Error performing operation")
		fmt.Println(components.Red.Render("\n  Could not create the handler file. Please check your permissions and disk space."))
		os.Exit(0)
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, map[string]string{
		"HandlerName": handlerName,
	})
	if err != nil {
		prompts.Error(components.Red.Render("Error performing operation"))
		fmt.Println(components.Red.Render("\n  Could not generate the handler file from the template. Please check the template syntax."))

		os.Exit(0)
		return err
	}

	prompts.Success("Done!!")

	fmt.Println(components.Green.Render(fmt.Sprintf("\n  Handler %s created successfully at: %s", handlerName, filePath)))

	return nil
}

func HandlerTargetDir(base, structure, service, domain, module, handlerName string) (string, error) {
	switch structure {
	case "clean_architecture":
		return filepath.Join(base, service, "application", "service"), nil
	case "domain_driven":
		domains, _ := os.ReadDir(base)
		var domNames []string
		for _, d := range domains {
			if d.IsDir() && d.Name() != "internal" && d.Name() != "modules" {
				domNames = append(domNames, d.Name())
			}
		}
		selectedDomain := domain
		return filepath.Join(base, service, selectedDomain, "handler"), nil
	case "modular":
		modules, _ := os.ReadDir(filepath.Join(base, "modules"))
		var modNames []string
		for _, m := range modules {
			if m.IsDir() {
				modNames = append(modNames, m.Name())
			}
		}
		selectedModule := module
		return filepath.Join(base, service, "modules", selectedModule, "handler"), nil
	case "layered":
		return filepath.Join(base, service, "handler"), nil
	case "custom":
		baseDir := filepath.Join(base, service)
		customDir := components.SelectCustomDirectory(baseDir, "handler")
		return customDir, nil
	default:
		return filepath.Join(base, service, "internal", "handler"), nil
	}
}
