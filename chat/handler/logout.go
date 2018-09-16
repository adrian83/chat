package handler

import (
	"net/http"

	"github.com/adrian83/chat/chat/session"
)

// LogoutHandler struct responsible for handling logout action.
type LogoutHandler struct {
	session *session.Session
}

// NewLogoutHandler returns new LogoutHandler struct.
func NewLogoutHandler(session *session.Session) *LogoutHandler {
	return &LogoutHandler{session: session}
}

// Logout processes logout action.
func (h *LogoutHandler) Logout(w http.ResponseWriter, req *http.Request) {

	if err := h.session.Remove(w, req); err != nil {
		RenderError500(w, err)
		return
	}

	http.Redirect(w, req, "/?reason=logout", http.StatusFound)
}
