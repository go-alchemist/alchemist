package templates

import (
	"text/template"
)

const serviceTemplate = `package service

// {{ .ServiceName }} provides business logic for the {{ .ServiceName }} entity.
type {{ .ServiceName }} struct {
	// TODO: Add service dependencies here
}

// New{{ .ServiceName }} creates a new {{ .ServiceName }} service.
func New{{ .ServiceName }}() *{{ .ServiceName }} {
	return &{{ .ServiceName }}{}
}
`

func GetServiceTemplate() (*template.Template, error) {
	return template.New("service").Parse(serviceTemplate)
}
