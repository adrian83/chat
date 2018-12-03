package handler

import (
	"fmt"
	"net/http"

	"github.com/adrian83/chat/chat/db"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler struct responsible for handling actions
// made on register html page.
type RegisterHandler struct {
	userRepo  *db.UserRepository
	templates *TemplateRepository
}

// NewRegisterHandler returns new RegisterHandler struct.
func NewRegisterHandler(templates *TemplateRepository, userRepo *db.UserRepository) *RegisterHandler {
	return &RegisterHandler{
		userRepo:  userRepo,
		templates: templates,
	}
}

// ShowRegisterPage renders register html page.
func (h *RegisterHandler) ShowRegisterPage(w http.ResponseWriter, req *http.Request) {
	RenderTemplate(w, h.templates.Register())
}

func (h *RegisterHandler) validateForm(username, password1, password2 string) []string {
	errors := make([]string, 0)

	if l := len(username); l < 3 || l > 200 {
		errors = append(errors, "Username should have more than 3 and less than 200 characters")
	}

	if l := len(password1); l < 3 || l > 200 {
		errors = append(errors, "Password should have more than 3 and less than 200 characters")
	}

	if l := len(password2); l < 3 || l > 200 {
		errors = append(errors, "Repeated password should have more than 3 and less than 200 characters")
	}

	if password1 != password2 {
		errors = append(errors, "Passwords should be the same")
	}

	return errors
}

// RegisterUser processes user registration form.
func (h *RegisterHandler) RegisterUser(w http.ResponseWriter, req *http.Request) {

	model := NewModel()

	if err := req.ParseForm(); err != nil {
		model.AddError(fmt.Sprintf("Cannot parse form: %v", err))
		RenderTemplateWithModel(w, h.templates.ServerError(), model)
		return
	}

	username := req.FormValue("username")
	password1 := req.FormValue("password1")
	password2 := req.FormValue("password2")

	validationErrors := h.validateForm(username, password1, password2)

	model["username"] = username
	model.AddErrors(validationErrors...)

	if model.HasErrors() {
		RenderTemplateWithModel(w, h.templates.Register(), model)
		return
	}

	user, err := h.userRepo.FindUser(username)
	if err != nil && err != db.ErrNotFound {
		model.AddError(fmt.Sprintf("Cannot get data about user: %v", err))
		RenderTemplateWithModel(w, h.templates.Login(), model)
		return
	}

	if !user.Empty() {
		model.AddError("User with this username already exists")
		RenderTemplateWithModel(w, h.templates.Register(), model)
		return
	}

	passBytes, err := bcrypt.GenerateFromPassword([]byte(password1), 1)
	if err != nil {
		model.AddError(fmt.Sprintf("Password encription failed: %v", err))
		RenderTemplateWithModel(w, h.templates.Login(), model)
		return
	}

	logger.Infof("Registering user with username: %s", username)

	if model.HasErrors() {
		RenderTemplateWithModel(w, h.templates.Register(), model)
		return
	}

	user = db.User{Login: username, Password: string(passBytes)}
	if err = h.userRepo.SaveUser(user); err != nil {
		model.AddError(fmt.Sprintf("Cannot store user data: %v", err))
		RenderTemplateWithModel(w, h.templates.Login(), model)
		return
	}

	http.Redirect(w, req, "/login", http.StatusFound)
}
