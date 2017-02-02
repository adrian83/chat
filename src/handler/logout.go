package handler

import (
	"net/http"
	mysession "session"
)

// LogoutHandler struct responsible for handling logout action.
type LogoutHandler struct {
	session *mysession.Session
}

// NewLogoutHandler returns new LogoutHandler struct.
func NewLogoutHandler(session *mysession.Session) *LogoutHandler {
	return &LogoutHandler{session: session}
}

// Logout processes logout action.
func (h *LogoutHandler) Logout(w http.ResponseWriter, req *http.Request) {

	if err := h.session.Remove(w, req); err != nil {
		RenderError500(w, err)
		return
	}

	http.Redirect(w, req, "/login", http.StatusFound)
}
