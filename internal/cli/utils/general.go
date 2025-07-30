package utils

import (
	"fmt"
	"os"

	"github.com/orochaa/go-clack/prompts"

	"github.com/go-alchemist/alchemist/internal/cli/components"
)

func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func PrintErrorAndReturn(msg string) error {
	prompts.Error(components.Red.Render("Error performing operation"))
	prompts.Outro(components.Red.Render(msg))
	return nil
}

func PrintSuccess(msg string) {
	prompts.Outro(components.Green.Render(msg))
}

func PrintSuccessf(format string, args ...interface{}) {
	fmt.Printf(components.Green.Render(format), args...)
}
