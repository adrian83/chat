package handler

import (
	"net/http"

	"github.com/adrian83/chat/chat/db"
	session "github.com/adrian83/go-redis-session"
	logger "github.com/sirupsen/logrus"
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

	sessionCookie, err := req.Cookie(sessionIDName)
	if err != nil {
		logger.Info("1")
		RenderTemplateWithModel(w, indexTmpl, model)
		return
	}

	session, err := h.sessionStore.Find(sessionCookie.Value)
	if err != nil {
		logger.Info("2")
		RenderError500(w, err)
		return
	}

	user := new(db.User)
	err = session.Get("user", user)
	logger.Infof("values %v", session.Values())
	if err != nil {
		logger.Info("3")
		RenderError500(w, err)
		return
	}

	model["loggedIn"] = !user.Empty()

	RenderTemplateWithModel(w, indexTmpl, model)
}
