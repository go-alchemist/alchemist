package make

import (
	"github.com/spf13/cobra"
)

var MakeCmd = &cobra.Command{
	Use:   "make",
	Short: "Scaffolding commands",
	Long:  "Generate files like models, migrations, handlers, etc.",
}

func init() {
	MakeCmd.AddCommand(MigrationCmd)
	MakeCmd.AddCommand(ModelCmd)
	MakeCmd.AddCommand(HandlerCmd)

}
