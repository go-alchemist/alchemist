package templates

import (
	"text/template"
)

const modelTemplate = `package models
type {{ .ModelName }} struct {
	// Fields here
}
`

func GetModelTemplate() (*template.Template, error) {
	return template.New("handler").Parse(modelTemplate)
}
