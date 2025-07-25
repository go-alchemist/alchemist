package make

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/go-alchemist/alchemist/internal/cli/response"
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
	content := fmt.Sprintf("package models\n\ntype %s struct {\n    // Fields here\n}\n", modelName)
	dir, _ := cmd.Flags().GetString("dir")
	if dir == "" {
		dir = "internal/models"
	}
	filePath := fmt.Sprintf("%s/%s.go", dir, strings.ToLower(modelName))
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		response.Error("Error creating model file: " + err.Error())
		return
	}

	response.Success(fmt.Sprintf("Model created: %s", filePath))
}
