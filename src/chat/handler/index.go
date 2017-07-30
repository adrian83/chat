package handler

import (
	"chat/session"
	"net/http"
)

var (
	indexTmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("index").WithTags("footer", "navigation", "head", "info").Build()
)

// IndexHandler struct responsible for handling actions
// made on index html page.
type IndexHandler struct {
	session *session.Session
}

// NewIndexHandler returns new IndexHandler struct.
func NewIndexHandler(session *session.Session) *IndexHandler {
	return &IndexHandler{session: session}
}

// ShowIndexPage renders Index page.
func (h *IndexHandler) ShowIndexPage(w http.ResponseWriter, req *http.Request) {
	model := NewModel()
	model["loggedIn"] = false

	if reason := req.URL.Query().Get("reason"); reason == "logout" {
		model.AddInfo("You have been logged out.")
	}

	sessionID := session.FindSessionID(req)
	if sessionID == "" {
		RenderTemplateWithModel(w, indexTmpl, model)
		return
	}

	user, err := h.session.FindUserData(sessionID)
	if err != nil {
		RenderError500(w, err)
		return
	}

	model["loggedIn"] = !user.Empty()

	RenderTemplateWithModel(w, indexTmpl, model)
}
