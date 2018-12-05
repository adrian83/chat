package handler

import (
	"html/template"
	"testing"

	"github.com/adrian83/chat/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestTemplatesShouldNotBeNilAndHaveProperNames(t *testing.T) {
	config := config.StaticsConfig{
		Path: "../../static",
	}

	templates := NewTemplateRepository(config)

	for _, tmpl := range []*template.Template{templates.Conversation, templates.Index, templates.Login, templates.Register, templates.ServerError} {
		assert.NotNil(t, tmpl)
		assert.Equal(t, "main.html", tmpl.Name(), "different name")
	}
}
