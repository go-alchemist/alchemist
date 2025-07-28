package setup

import "github.com/orochaa/go-clack/prompts"

func (m *model) DefaultSettings() {
	confirmed, err := prompts.Confirm(prompts.ConfirmParams{
		Message:      "Do you want to use the default settings?",
		InitialValue: true,
	})

	prompts.ExitOnError(err)

	m.Default = confirmed
}
