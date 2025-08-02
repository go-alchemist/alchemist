package templates

import (
	"text/template"
)

const usecaseTemplate = `package usecase

// {{ .UsecaseName }}Usecase defines the application logic for {{ .UsecaseName }}.
type {{ .UsecaseName }}Usecase struct {
	// TODO: Add usecase dependencies here
}

// New{{ .UsecaseName }}Usecase creates a new {{ .UsecaseName }}Usecase.
func New{{ .UsecaseName }}Usecase() *{{ .UsecaseName }}Usecase {
	return &{{ .UsecaseName }}Usecase{}
}
`

func GetUsecaseTemplate() (*template.Template, error) {
	return template.New("usecase").Parse(usecaseTemplate)
}
