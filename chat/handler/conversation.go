package handler

import (
	"net/http"

	"github.com/adrian83/chat/chat/db"
	session "github.com/adrian83/go-redis-session"
	logger "github.com/sirupsen/logrus"
)

var (
	conversationTmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("conversation").WithTags("footer", "navigation", "head").Build()
)

// ConversationHandler struct responsible for handling actions
// made on index html page.
type ConversationHandler struct {
	sessionStore session.Store
}

// NewConversationHandler returns new ConversationHandler struct.
func NewConversationHandler(sessionStore session.Store) *ConversationHandler {
	return &ConversationHandler{sessionStore: sessionStore}
}

// ShowConversationPage renders conversation page.
func (h *ConversationHandler) ShowConversationPage(w http.ResponseWriter, req *http.Request) {

	sessionCookie, err := req.Cookie(sessionIDName)
	if err != nil {
		logger.Info("1")
		RenderError500(w, err)
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
	if err != nil {
		logger.Info("3")
		RenderError500(w, err)
		return
	}

	if user.Empty() {
		logger.Info("4")
		http.Redirect(w, req, "/login", http.StatusFound)
		return
	}

	RenderTemplateWithModel(w, conversationTmpl, Model(map[string]interface{}{"sessionId": sessionCookie.Value}))

}
