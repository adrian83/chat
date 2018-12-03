package handler

import (
	"net/http"
)

var (
	error500Tmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("error500").WithTags("footer", "navigation", "head").Build()
)

// RenderError500 renders page with error message
func RenderError500(w http.ResponseWriter, err error) {
	model := NewModel()
	model["message"] = err.Error()

	RenderTemplateWithModel(w, error500Tmpl, model)
}
