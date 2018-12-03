package handler

import (
	"fmt"
	"net/http"

	"github.com/adrian83/chat/chat/db"
	session "github.com/adrian83/go-redis-session"
)

var (
	indexTmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("index").WithTags("footer", "navigation", "head", "info").Build()
)

// IndexHandler struct responsible for handling actions
// made on index html page.
type IndexHandler struct {
	sessionStore session.Store
}

// NewIndexHandler returns new IndexHandler struct.
func NewIndexHandler(sessionStore session.Store) *IndexHandler {
	return &IndexHandler{sessionStore: sessionStore}
}

// ShowIndexPage renders Index page.
func (h *IndexHandler) ShowIndexPage(w http.ResponseWriter, req *http.Request) {
	model := NewModel()
	model["loggedIn"] = false

	if reason := req.URL.Query().Get("reason"); reason == "logout" {
		model.AddInfo("You have been logged out.")
	}

	sessionID, err := ReadSessionIDFromCookie(req)
	if err != nil {
		model.AddError(fmt.Sprintf("Cannot find session id: %v", err))
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	userSession, err := h.sessionStore.Find(sessionID)
	if err != nil {
		model.AddError(fmt.Sprintf("Cannot find user session: %v", err))
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	user := new(db.User)
	if err = userSession.Get("user", user); err != nil {
		model.AddError(fmt.Sprintf("Cannot get data about user: %v", err))
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	model["loggedIn"] = !user.Empty()

	RenderTemplateWithModel(w, indexTmpl, model)
}
