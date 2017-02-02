package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

// Model is a struct with all needed data for rendering template.
type Model map[string]interface{}

// NewModel function returns new instance of Model struct.
func NewModel() Model {
	e := make([]string, 0)
	m := Model(make(map[string]interface{}))
	m["errors"] = e
	return m
}

// AddError adds error message to model.
func (m Model) AddError(errMsg string) {
	e := m["errors"].([]string)
	m["errors"] = append(e, errMsg)
}

// HasErrors returns 'true' if at least one erros message is in the
// model, 'false' otherwise.
func (m Model) HasErrors() bool {
	e := m["errors"].([]string)
	return len(e) > 0
}

// RenderTemplate renders given template with empty model.
func RenderTemplate(w http.ResponseWriter, tmpl *template.Template) bool {
	return RenderTemplateWithModel(w, tmpl, NewModel())
}

// RenderTemplateWithModel renders given template with given model.
func RenderTemplateWithModel(w http.ResponseWriter, tmpl *template.Template, model Model) bool {
	if err := tmpl.ExecuteTemplate(w, "base", model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	return true
}

// NewTemplateBuilder returns new instance of TemplateBuilder struct.
func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{paths: make([]string, 0)}
}

// TemplateBuilder is a struct used for creating
// instance of html/template.Template.
type TemplateBuilder struct {
	paths []string
}

// WithMainTemplate sets template which will be used to create Template struct.
func (b *TemplateBuilder) WithMainTemplate(name string) *TemplateBuilder {
	b.paths = append(b.paths, fmt.Sprintf("static/pages/templates/%s.html", name))
	return b
}

// WithContent sets content which will be used to create Template struct.
func (b *TemplateBuilder) WithContent(name string) *TemplateBuilder {
	b.paths = append(b.paths, fmt.Sprintf("static/pages/%s.html", name))
	return b
}

// WithTags sets tags (footer, header...) which will be used
// to create Template struct.
func (b *TemplateBuilder) WithTags(names ...string) *TemplateBuilder {
	for _, name := range names {
		b.paths = append(b.paths, fmt.Sprintf("static/pages/fragments/%s.html", name))
	}
	return b
}

// Build creates and returns previously prepared Template struct.
func (b *TemplateBuilder) Build() *template.Template {
	return template.Must(template.ParseFiles(b.paths...))
}
