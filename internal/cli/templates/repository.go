package templates

import (
	"text/template"
)

const repositoryTemplate = `package repository

// {{ .RepositoryName }}Repository defines the interface for storing {{ .RepositoryName }} entities.
type {{ .RepositoryName }}Repository interface {
	// TODO: Add repository methods (e.g., Save, Find, Delete)
}
`

func GetRepositoryTemplate() (*template.Template, error) {
	return template.New("repository").Parse(repositoryTemplate)
}
