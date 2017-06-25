package handler

import (
	"net/http"

	"chat/db"
	"chat/logger"

	"golang.org/x/crypto/bcrypt"
)

var (
	registerTmpl = NewTemplateBuilder().WithMainTemplate("main").WithContent("register").WithTags("footer", "navigation", "head", "errors").Build()
)

// RegisterHandler struct responsible for handling actions
// made on register html page.
type RegisterHandler struct {
	userRepo *db.UserRepository
}

// NewRegisterHandler returns new RegisterHandler struct.
func NewRegisterHandler(userRepo *db.UserRepository) *RegisterHandler {
	return &RegisterHandler{userRepo: userRepo}
}

// ShowRegisterPage renders register html page.
func (h *RegisterHandler) ShowRegisterPage(w http.ResponseWriter, req *http.Request) {
	RenderTemplate(w, registerTmpl)
}

func (h *RegisterHandler) validateForm(username, password1, password2 string, model Model) {
	if l := len(username); l < 3 || l > 200 {
		model.AddError("Username should have more than 3 and less than 200 characters")
	}

	if l := len(password1); l < 3 || l > 200 {
		model.AddError("Password should have more than 3 and less than 200 characters")
	}

	if l := len(password2); l < 3 || l > 200 {
		model.AddError("Repeated password should have more than 3 and less than 200 characters")
	}

	if password1 != password2 {
		model.AddError("Passwords should be the same")
	}
}

// RegisterUser processes user registration form.
func (h *RegisterHandler) RegisterUser(w http.ResponseWriter, req *http.Request) {

	if err := req.ParseForm(); err != nil {
		RenderError500(w, err)
		return
	}

	username := req.FormValue("username")
	password1 := req.FormValue("password1")
	password2 := req.FormValue("password2")

	model := NewModel()
	model["username"] = username

	h.validateForm(username, password1, password2, model)

	if model.HasErrors() {
		RenderTemplateWithModel(w, registerTmpl, model)
		return
	}

	logger.Info("RegisterHandler", "RegisterUser", "Basic validation passed")

	_, exists, err := h.userRepo.FindUser(username)
	if err != nil {
		RenderError500(w, err)
		return
	}

	if exists {
		model.AddError("User with this username already exists")
		RenderTemplateWithModel(w, registerTmpl, model)
		return
	}

	passBytes, err := bcrypt.GenerateFromPassword([]byte(password1), 1)
	if err != nil {
		RenderError500(w, err)
		return
	}

	logger.Infof("RegisterHandler", "RegisterUser", "Registering user with username: %s", username)

	if model.HasErrors() {
		RenderTemplateWithModel(w, registerTmpl, model)
		return
	}

	user := db.User{Name: username, Password: string(passBytes)}
	if err = h.userRepo.SaveUser(user); err != nil {
		RenderError500(w, err)
		return
	}

	http.Redirect(w, req, "/login", http.StatusFound)
}
