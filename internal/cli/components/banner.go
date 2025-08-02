package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func Banner() string {
	text := "Alchemist - CLI to projects in Go"

	colors := []string{
		"#42d392", "#3fd399", "#3bd3a1", "#37d3a8", "#33d3af", "#30d3b7",
		"#2cd3be", "#28d3c5", "#24d3cd", "#20d3d4", "#1ed0d9", "#1dbad9",
		"#1ba4d9", "#198ed9", "#1878d9", "#1662d9", "#154cd9", "#1336d9",
	}

	if termenv.EnvColorProfile() == termenv.TrueColor {
		var b strings.Builder
		for i, r := range text {
			col := colors[i%len(colors)]
			s := lipgloss.NewStyle().Foreground(lipgloss.Color(col)).Render(string(r))
			b.WriteString(s)
		}
		return b.String()
	}

	green := lipgloss.Color("#42d392")
	return lipgloss.NewStyle().
		Foreground(green).
		Bold(true).
		Render(text)
}
