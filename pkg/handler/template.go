package handler

import (
	"html/template"
)

// NewTemplateRepository returns new TemplateRepository.
func NewTemplateRepository(templatesPath string) *TemplateRepository {
	return &TemplateRepository{
		Login:        NewTemplateBuilder(templatesPath).WithTemplate("main").WithContent("login").WithTags("footer", "navigation", "head", "errors").Build(),
		Conversation: NewTemplateBuilder(templatesPath).WithTemplate("main").WithContent("conversation").WithTags("errors", "footer", "navigation", "head").Build(),
		ServerError:  NewTemplateBuilder(templatesPath).WithTemplate("main").WithContent("error500").WithTags("footer", "errors", "navigation", "head").Build(),
		Index:        NewTemplateBuilder(templatesPath).WithTemplate("main").WithContent("index").WithTags("errors", "footer", "navigation", "head", "info").Build(),
		Register:     NewTemplateBuilder(templatesPath).WithTemplate("main").WithContent("register").WithTags("footer", "navigation", "head", "errors").Build(),
	}
}

// TemplateRepository keeps all templates used in application.
type TemplateRepository struct {
	Login        *template.Template
	Conversation *template.Template
	ServerError  *template.Template
	Index        *template.Template
	Register     *template.Template
}
