package handler

import (
	"html/template"

	"github.com/adrian83/chat/chat/config"
)

func NewTemplateRepository(statics config.StaticsConfig) *TemplateRepository {
	return &TemplateRepository{
		statics: statics,
	}
}

type TemplateRepository struct {
	statics      config.StaticsConfig
	login        *template.Template
	conversation *template.Template
	serverError  *template.Template
	index        *template.Template
	register     *template.Template
}

func (tr *TemplateRepository) Login() *template.Template {
	if tr.login == nil {
		tr.login = NewTemplateBuilder(tr.statics.Path).WithMainTemplate("main").WithContent("login").WithTags("footer", "navigation", "head", "errors").Build()
	}
	return tr.login
}

func (tr *TemplateRepository) Conversation() *template.Template {
	if tr.conversation == nil {
		tr.conversation = NewTemplateBuilder(tr.statics.Path).WithMainTemplate("main").WithContent("conversation").WithTags("footer", "navigation", "head").Build()
	}
	return tr.conversation
}

func (tr *TemplateRepository) ServerError() *template.Template {
	if tr.serverError == nil {
		tr.serverError = NewTemplateBuilder(tr.statics.Path).WithMainTemplate("main").WithContent("error500").WithTags("footer", "errors", "navigation", "head").Build()
	}
	return tr.serverError
}

func (tr *TemplateRepository) Index() *template.Template {
	if tr.index == nil {
		tr.index = NewTemplateBuilder(tr.statics.Path).WithMainTemplate("main").WithContent("index").WithTags("footer", "navigation", "head", "info").Build()
	}
	return tr.index
}

func (tr *TemplateRepository) Register() *template.Template {
	if tr.register == nil {
		tr.register = NewTemplateBuilder(tr.statics.Path).WithMainTemplate("main").WithContent("register").WithTags("footer", "navigation", "head", "errors").Build()
	}
	return tr.register
}
