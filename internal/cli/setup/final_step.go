package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/orochaa/go-clack/prompts"

	"github.com/go-alchemist/alchemist/internal/cli/components"
)

func (m *model) FinalSteps(features map[string]bool) {
	prompts.Outro("Done. Now run:")

	marginCmd := lipgloss.NewStyle().MarginLeft(2)

	block := ""
	block += marginCmd.Render(components.Green.Render(fmt.Sprintf("cd %s", m.ProjectName))) + "\n"
	block += marginCmd.Render(components.Green.Render("go mod tidy")) + "\n"

	if features["microservice_architecture"] {
		block += marginCmd.Render(components.Green.Render(`for d in ./modules/* ; do (cd "$d" && go mod tidy); done`)) + "\n"
	}

	fmt.Print("\n" + block)

	if _, err := os.Stat(filepath.Join(m.ProjectName, ".git")); os.IsNotExist(err) {
		fmt.Println()
		fmt.Print(components.Dim.Render("|"))
		fmt.Println(" Optional: Initialize Git in your project directory with:")
		fmt.Println()
		fmt.Println(marginCmd.Render(components.Green.Render(` git init && git add -A && git commit -m "initial commit"`)))
	}
	fmt.Println()
	fmt.Print(components.Dim.Render("|"))
	fmt.Println(components.Yellow.Render(" There are additional configurations that can be made in the .yaml file. Read the documentation for more details."))
}
