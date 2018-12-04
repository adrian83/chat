package handler

import (
	"fmt"
	"net/http"

	"github.com/adrian83/chat/pkg/db"
	session "github.com/adrian83/go-redis-session"
)

// IndexHandler struct responsible for handling actions
// made on index html page.
type IndexHandler struct {
	sessionStore session.Store
	templates    *TemplateRepository
}

// NewIndexHandler returns new IndexHandler struct.
func NewIndexHandler(templates *TemplateRepository, sessionStore session.Store) *IndexHandler {
	return &IndexHandler{
		sessionStore: sessionStore,
		templates:    templates,
	}
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
		RenderTemplateWithModel(w, h.templates.Login, model)
		return
	}

	userSession, err := h.sessionStore.Find(sessionID)
	if err != nil {
		model.AddError(fmt.Sprintf("Cannot find user session: %v", err))
		RenderTemplateWithModel(w, h.templates.Login, model)
		return
	}

	user := new(db.User)
	if err = userSession.Get("user", user); err != nil {
		model.AddError(fmt.Sprintf("Cannot get data about user: %v", err))
		RenderTemplateWithModel(w, h.templates.Login, model)
		return
	}

	model["loggedIn"] = !user.Empty()

	RenderTemplateWithModel(w, h.templates.Index, model)
}
