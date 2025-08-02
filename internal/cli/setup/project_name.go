package setup

import (
	"errors"
	"os"

	"github.com/orochaa/go-clack/prompts"
)

func (m *model) SelectProjectName() {
	name, err := prompts.Text(prompts.TextParams{
		Message:  "Project name (target directory)?:",
		Required: true,
		Validate: func(value string) error {
			if value == "" {
				return errors.New("Project name cannot be empty")
			}
			if _, err := os.Stat(value); !os.IsNotExist(err) {
				return errors.New("Project name already exists, please choose a different name")
			}
			return nil
		},
	})
	prompts.ExitOnError(err)

	m.ProjectName = name
}
