package handler

import (
	"net/http"

	"github.com/adrian83/go-redis-session"
)

// LogoutHandler struct responsible for handling logout action.
type LogoutHandler struct {
	sessionStore session.Store
}

// NewLogoutHandler returns new LogoutHandler struct.
func NewLogoutHandler(sessionStore session.Store) *LogoutHandler {
	return &LogoutHandler{
		sessionStore: sessionStore,
	}
}

// Logout processes logout action.
func (h *LogoutHandler) Logout(w http.ResponseWriter, req *http.Request) {

	sessionCookie, err := req.Cookie(sessionIDName)
	if err != nil {
		RenderError500(w, err)
		return
	}

	if err := h.sessionStore.Delete(sessionCookie.Value); err != nil {
		RenderError500(w, err)
		return
	}

	http.Redirect(w, req, "/?reason=logout", http.StatusFound)
}
