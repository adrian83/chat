package handler

import (
	"fmt"

	"github.com/adrian83/chat/chat/db"

	"net/http"

	"github.com/adrian83/go-redis-session"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	loginTmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("login").WithTags("footer", "navigation", "head", "errors").Build()
)

// LoginHandler struct responsible for handling actions
// made on login html page.
type LoginHandler struct {
	userRepo     *db.UserRepository
	sessionStore session.Store
}

// NewLoginHandler returns new LoginHandler struct.
func NewLoginHandler(userRepo *db.UserRepository, sessionStore session.Store) *LoginHandler {
	return &LoginHandler{
		userRepo:     userRepo,
		sessionStore: sessionStore,
	}
}

// ShowLoginPage renders login html page.
func (h *LoginHandler) ShowLoginPage(w http.ResponseWriter, req *http.Request) {
	RenderTemplate(w, loginTmpl)
}

// LoginUser processes user login form.
func (h *LoginHandler) LoginUser(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		RenderError500(w, err)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	model := NewModel()
	model["username"] = username

	if username == "" {
		model.AddError("Username cannot be empty")
	}

	if password == "" {
		model.AddError("Password cannot be empty")
	}

	if model.HasErrors() {
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	user, err := h.userRepo.FindUser(username)
	if err != nil {
		model.AddError(fmt.Sprintf("Cannot get data about user: %v", err))
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	if user.Empty() {
		model.AddError("User with this username doesn't exist")
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		model.AddError(fmt.Sprintf("Passwords don't match: %v", err))
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	if err = h.storeInSession(user, w); err != nil {
		model.AddError(fmt.Sprintf("Cannot create session: %v", err))
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	http.Redirect(w, req, "/conversation", http.StatusFound)

}

func (h *LoginHandler) storeInSession(user db.User, w http.ResponseWriter) error {
	sessionID := uuid.New().String()
	session, err := h.sessionStore.Create(sessionID, defSessionDuration)
	if err != nil {
		return err
	}

	if err = session.Add("user", user); err != nil {
		return err
	}

	if err = h.sessionStore.Save(session); err != nil {
		return err
	}

	storeSessionCookie(sessionID, w)

	return nil
}
