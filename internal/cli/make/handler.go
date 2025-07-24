package make

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/go-alchemist/alchemist/internal/cli/response"
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
	content := fmt.Sprintf(`package handlers
import (
    "net/http"
)

// %s handles HTTP requests.
func %s(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement handler logic
}
`, handlerName, handlerName)
	dir, _ := cmd.Flags().GetString("dir")
	if dir == "" {
		dir = "internal/handlers"
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	filePath := fmt.Sprintf("%s/%s.go", dir, strings.ToLower(handlerName))
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		response.Error(fmt.Errorf("error creating handler: %v", err))
		return
	}

	response.Success(fmt.Sprintf("Handler created: %s", filePath))
}
