package handler

import (
	"net/http"

	websession "session"
)

var (
	conversationTmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("conversation").WithTags("footer", "navigation", "head").Build()
)

// ConversationHandler struct responsible for handling actions
// made on index html page.
type ConversationHandler struct {
	session *websession.Session
}

// NewConversationHandler returns new ConversationHandler struct.
func NewConversationHandler(session *websession.Session) *ConversationHandler {
	return &ConversationHandler{session: session}
}

// ShowConversationPage renders conversation page.
func (h *ConversationHandler) ShowConversationPage(w http.ResponseWriter, req *http.Request) {

	sessionID := websession.FindSessionID(req)
	if sessionID == "" {
		http.Redirect(w, req, "/login", http.StatusFound)
		return
	}

	user, err := h.session.FindUserData(sessionID)
	if err != nil {
		RenderError500(w, err)
		return
	}

	if user.Empty() {
		http.Redirect(w, req, "/login", http.StatusFound)
		return
	}

	RenderTemplateWithModel(w, conversationTmpl, Model(map[string]interface{}{"sessionId": sessionID}))

}
