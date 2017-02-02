package handler

import (
	"net/http"

	websession "session"
)

var (
	talkTmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("index").WithTags("footer", "navigation", "head").Build()
)

// IndexHandler struct responsible for handling actions
// made on index html page.
type IndexHandler struct {
	session *websession.Session
}

// NewIndexHandler returns new IndexHandler struct.
func NewIndexHandler(session *websession.Session) *IndexHandler {
	return &IndexHandler{session: session}
}

// ShowIndexPage renders Index page.
func (h *IndexHandler) ShowIndexPage(w http.ResponseWriter, req *http.Request) {

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

	RenderTemplateWithModel(w, talkTmpl, Model(map[string]interface{}{"sessionId": sessionID}))

}
