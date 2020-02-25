package handler

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplatesShouldNotBeNilAndHaveProperNames(t *testing.T) {
	// given
	staticsPath := "../../static"

	// when
	templates := NewTemplateRepository(staticsPath)

	// then
	for _, tmpl := range []*template.Template{templates.Conversation, templates.Index, templates.Login, templates.Register, templates.ServerError} {
		assert.NotNil(t, tmpl)
		assert.Equal(t, "main.html", tmpl.Name(), "different name")
	}
}
