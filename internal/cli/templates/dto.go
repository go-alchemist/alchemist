package templates

import (
	"text/template"
)

const dtoTemplate = `package dto

// {{ .DTOName }}DTO is a Data Transfer Object for {{ .DTOName }}.
type {{ .DTOName }}DTO struct {
	// TODO: Add DTO fields here
}
`

func GetDTOTemplate() (*template.Template, error) {
	return template.New("dto").Parse(dtoTemplate)
}
