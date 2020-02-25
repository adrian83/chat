package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/adrian83/chat/pkg/user"
)

const (
	errorsField = "errors"
	infoField   = "info"
	userField   = "user"
)

// Model is a struct with all needed data for rendering template.
type Model map[string]interface{}

// NewModel function returns new instance of Model struct.
func NewModel() Model {
	m := Model(make(map[string]interface{}))
	m[errorsField] = make([]string, 0)
	m[infoField] = make([]string, 0)
	return m
}

func (m Model) AddUser(usr *user.User) {
	if usr != nil && !usr.Empty() {
		m[userField] = usr
	}
}

// AddError adds error message to model.
func (m Model) AddError(errMsg string) {
	e := m[errorsField].([]string)
	m[errorsField] = append(e, errMsg)
}

// AddErrors adds multiple error messages to model.
func (m Model) AddErrors(errMsgs ...error) {
	e := m[errorsField].([]string)
	for _, err := range errMsgs {
		e = append(e, err.Error())
	}
	m[errorsField] = e
}

// AddInfo adds info message to model.
func (m Model) AddInfo(infoMsg string) {
	e := m[infoField].([]string)
	m[infoField] = append(e, infoMsg)
}

// HasErrors returns 'true' if at least one erros message is in the
// model, 'false' otherwise.
func (m Model) HasErrors() bool {
	e := m[errorsField].([]string)
	return len(e) > 0
}

// RenderTemplate renders given template with empty model.
func RenderTemplate(w http.ResponseWriter, tmpl *template.Template) {
	RenderTemplateWithModel(w, tmpl, NewModel())
}

// RenderTemplateWithModel renders given template with given model.
func RenderTemplateWithModel(w http.ResponseWriter, tmpl *template.Template, model Model) {
	if err := tmpl.ExecuteTemplate(w, "base", model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NewTemplateBuilder returns new instance of TemplateBuilder struct.
func NewTemplateBuilder(path string) *TemplateBuilder {
	return &TemplateBuilder{
		base:  path,
		paths: make([]string, 0),
	}
}

// TemplateBuilder is a struct used for creating
// instance of html/template.Template.
type TemplateBuilder struct {
	paths []string
	base  string
}

// WithMainTemplate sets template which will be used to create Template struct.
func (b *TemplateBuilder) WithTemplate(name string) *TemplateBuilder {
	b.paths = append(b.paths, fmt.Sprintf("%v/pages/templates/%s.html", b.base, name))
	return b
}

// WithContent sets content which will be used to create Template struct.
func (b *TemplateBuilder) WithContent(name string) *TemplateBuilder {
	b.paths = append(b.paths, fmt.Sprintf("%v/pages/%s.html", b.base, name))
	return b
}

// WithTags sets tags (footer, header...) which will be used
// to create Template struct.
func (b *TemplateBuilder) WithTags(names ...string) *TemplateBuilder {
	for _, name := range names {
		b.paths = append(b.paths, fmt.Sprintf("%v/pages/fragments/%s.html", b.base, name))
	}
	return b
}

// Build creates and returns previously prepared Template struct.
func (b *TemplateBuilder) Build() *template.Template {
	return template.Must(template.ParseFiles(b.paths...))
}
