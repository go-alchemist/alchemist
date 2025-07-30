package templates

import (
	"text/template"
)

const modelTemplate = `package model

// {{ .ModelName }} represents the {{ .ModelName }} entity.
type {{ .ModelName }} struct {
	// TODO: Add fields
}
`

func GetModelTemplate() (*template.Template, error) {
	return template.New("model").Parse(modelTemplate)
}
