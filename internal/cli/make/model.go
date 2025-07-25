package make

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/go-alchemist/alchemist/internal/cli/config"
	"github.com/go-alchemist/alchemist/internal/cli/response"
	"github.com/go-alchemist/alchemist/internal/cli/templates"
)

var ModelCmd = &cobra.Command{
	Use:   "model [name]",
	Short: "Create a new model",
	Args:  cobra.ExactArgs(1),
	Run:   MakeModel,
}

func init() {
	ModelCmd.Flags().String("dir", "internal/models", "Directory for models")
}

func MakeModel(cmd *cobra.Command, args []string) {
	modelName := args[0]

	flagDir, _ := cmd.Flags().GetString("dir")
	originalDir := flagDir
	if originalDir == "" {
		originalDir = config.Config.GetString("paths.models")
		if originalDir == "" {
			originalDir = "internal/models"
		}
	}

	dir := response.SelectMicroserviceIfEnabled()
	if dir == "" {
		response.Error("Microservice feature is not enabled or directory not found")
		return
	}

	dir = path.Join(dir, originalDir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	filePath := path.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(modelName)))

	tmpl, err := templates.GetModelTemplate()
	if err != nil {
		response.Error("Error loading model template: " + err.Error())
		return
	}

	f, err := os.Create(filePath)
	if err != nil {
		response.Error("Error creating model file: " + err.Error())
		return
	}
	defer f.Close()

	err = tmpl.Execute(f, map[string]string{
		"ModelName": modelName,
	})
	if err != nil {
		response.Error("Error executing model template: " + err.Error())
		return
	}

	response.Success(fmt.Sprintf("Model created: %s", filePath))
}
