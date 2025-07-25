package templates

import (
	"text/template"
)

const handlerTemplate = `package handlers
import (
    "net/http"
)

// {{ .HandlerName }} handles HTTP requests.
func {{ .HandlerName }}(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement handler logic
}
`

func GetHandlerTemplate() (*template.Template, error) {
	return template.New("handler").Parse(handlerTemplate)
}
