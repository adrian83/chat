package handler

import (
	"db"
	gsession "session"

	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var (
	loginTmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("login").WithTags("footer", "navigation", "head", "errors").Build()
)

// LoginHandler struct responsible for handling actions
// made on login html page.
type LoginHandler struct {
	userRepo *db.UserRepository
	session  *gsession.Session
}

// NewLoginHandler returns new LoginHandler struct.
func NewLoginHandler(userRepo *db.UserRepository, session *gsession.Session) *LoginHandler {
	return &LoginHandler{userRepo: userRepo, session: session}
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

	user, ok, err := h.userRepo.FindUser(username)
	if err != nil {
		RenderError500(w, err)
		return
	}

	if !ok {
		model.AddError("User with this username doesn't exist")
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		model.AddError("Passwords don't match: " + err.Error())
		RenderTemplateWithModel(w, loginTmpl, model)
		return
	}

	_, err2 := h.session.StoreUserData(w, user)
	if err2 != nil {
		RenderError500(w, err2)
		return
	}

	http.Redirect(w, req, "/", http.StatusFound)

}
