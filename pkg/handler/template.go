package handler

import (
	"html/template"

	"github.com/adrian83/chat/chat/config"
)

// NewTemplateRepository returns new TemplateRepository.
func NewTemplateRepository(statics config.StaticsConfig) *TemplateRepository {
	return &TemplateRepository{
		statics:      statics,
		Login:        NewTemplateBuilder(statics.Path).WithMainTemplate("main").WithContent("login").WithTags("footer", "navigation", "head", "errors").Build(),
		Conversation: NewTemplateBuilder(statics.Path).WithMainTemplate("main").WithContent("conversation").WithTags("errors", "footer", "navigation", "head").Build(),
		ServerError:  NewTemplateBuilder(statics.Path).WithMainTemplate("main").WithContent("error500").WithTags("footer", "errors", "navigation", "head").Build(),
		Index:        NewTemplateBuilder(statics.Path).WithMainTemplate("main").WithContent("index").WithTags("errors", "footer", "navigation", "head", "info").Build(),
		Register:     NewTemplateBuilder(statics.Path).WithMainTemplate("main").WithContent("register").WithTags("footer", "navigation", "head", "errors").Build(),
	}
}

// TemplateRepository keeps all templates used in application.
type TemplateRepository struct {
	statics      config.StaticsConfig
	Login        *template.Template
	Conversation *template.Template
	ServerError  *template.Template
	Index        *template.Template
	Register     *template.Template
}
