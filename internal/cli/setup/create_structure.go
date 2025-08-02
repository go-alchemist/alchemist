package setup

import "github.com/orochaa/go-clack/prompts"

func (m *model) CreateStructure() {
	confirmed, err := prompts.Confirm(prompts.ConfirmParams{
		Message:      "Do you want to create the folder structure?",
		InitialValue: true,
	})

	prompts.ExitOnError(err)

	m.CreateStructureFolder = confirmed
}
