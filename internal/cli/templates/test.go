package templates

import (
	"text/template"
)

const testTemplate = `package {{ .TestName | ToLower }}_test

import (
	"testing"
)

func Test{{ .TestName }}(t *testing.T) {
	// TODO: Write test cases for {{ .TestName }}
}
`

func GetTestTemplate() (*template.Template, error) {
	return template.New("test").Parse(testTemplate)
}
