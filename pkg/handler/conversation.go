package handler

import (
	"fmt"
	"net/http"

	"github.com/adrian83/chat/pkg/user"
	session "github.com/adrian83/go-redis-session"
)

// ConversationHandler struct responsible for handling actions
// made on index html page.
type ConversationHandler struct {
	sessionStore *session.Store
	templates    *TemplateRepository
}

// NewConversationHandler returns new ConversationHandler struct.
func NewConversationHandler(templates *TemplateRepository, sessionStore *session.Store) *ConversationHandler {
	return &ConversationHandler{
		sessionStore: sessionStore,
		templates:    templates,
	}
}

// ShowConversationPage renders conversation page.
func (h *ConversationHandler) ShowConversationPage(w http.ResponseWriter, req *http.Request) {
	model := NewModel()

	sessionCookie, err := req.Cookie(sessionIDName)
	if err != nil {
		model.AddError(fmt.Sprintf("Cannot find session cookie: %v", err))
		RenderTemplateWithModel(w, h.templates.ServerError, model)
		return
	}

	session, err := h.sessionStore.Find(sessionCookie.Value)
	if err != nil {
		model.AddError(fmt.Sprintf("Cannot find user session: %v", err))
		RenderTemplateWithModel(w, h.templates.ServerError, model)
		return
	}

	var user user.User
	if err = session.Get("user", &user); err != nil {
		model.AddError(fmt.Sprintf("Cannot find user data in session: %v", err))
		RenderTemplateWithModel(w, h.templates.ServerError, model)
		return
	}

	if user.Empty() {
		http.Redirect(w, req, "/login", http.StatusFound)
		return
	}

	var modelDict = map[string]interface{}{
		"sessionId": sessionCookie.Value,
		"username":  user.Name(),
	}

	RenderTemplateWithModel(w, h.templates.Conversation, Model(modelDict))
}
