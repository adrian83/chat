package handler

import (
	"net/http"

	websession "session"
)

var (
	indexTmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("index").WithTags("footer", "navigation", "head").Build()
)

// IndexHandler struct responsible for handling actions
// made on index html page.
type IndexHandler struct {
	session *websession.Session
}

// NewIndexHandler returns new IndexHandler struct.
func NewIndexHandler(session *websession.Session) *IndexHandler {
	return &IndexHandler{session: session}
}

// ShowIndexPage renders Index page.
func (h *IndexHandler) ShowIndexPage(w http.ResponseWriter, req *http.Request) {
	RenderTemplate(w, indexTmpl)
}
