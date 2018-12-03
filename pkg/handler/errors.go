package handler

import (
	"html/template"
	"net/http"
)

// RenderError500 renders page with error message
func RenderError500(w http.ResponseWriter, template *template.Template, err error) {
	model := NewModel()
	model["message"] = err.Error()

	RenderTemplateWithModel(w, template, model)
}
