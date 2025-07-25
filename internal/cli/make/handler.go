package make

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/go-alchemist/alchemist/internal/cli/config"
	"github.com/go-alchemist/alchemist/internal/cli/response"
	"github.com/go-alchemist/alchemist/internal/cli/templates"
)

var HandlerCmd = &cobra.Command{
	Use:   "handler [name]",
	Short: "Create a new handler",
	Args:  cobra.ExactArgs(1),
	Run:   MakeHandler,
}

func init() {
	HandlerCmd.Flags().String("dir", "internal/handlers", "Directory for handlers")
}

func MakeHandler(cmd *cobra.Command, args []string) {
	handlerName := args[0]

	originalDir := config.Config.GetString("paths.handlers")
	if originalDir == "" {
		originalDir = "internal/handlers"
	}

	dir := response.SelectMicroserviceIfEnabled(originalDir)
	if dir == "" {
		response.Error("Microservice feature is not enabled or directory not found")
		return
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	filePath := fmt.Sprintf("%s/%s.go", dir, strings.ToLower(handlerName))
	tmpl, err := templates.GetHandlerTemplate()
	if err != nil {
		response.Error("Error loading handler template")
		return
	}
	f, err := os.Create(filePath)
	if err != nil {
		response.Error("Error creating handler file: " + err.Error())
		return
	}
	defer f.Close()

	err = tmpl.Execute(f, map[string]string{
		"HandlerName": handlerName,
	})
	if err != nil {
		response.Error("Error executing handler template: " + err.Error())
		return
	}

	response.Success(fmt.Sprintf("Handler created: %s", filePath))
}
