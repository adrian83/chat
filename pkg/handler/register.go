package handler

import (
	"fmt"
	"net/http"

	"github.com/adrian83/chat/pkg/db"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	minUsernameLen = 3
	maxUsernameLen = 200
	minPasswordLen = 3
	maxPasswordLen = 200
)

var (
	errInvalidUsername    = fmt.Errorf("Username should have more than 3 and less than 200 characters")
	errInvalidPassword1   = fmt.Errorf("Password should have more than 3 and less than 200 characters")
	errInvalidPassword2   = fmt.Errorf("Repeated password should have more than 3 and less than 200 characters")
	errDifferendPasswords = fmt.Errorf("Passwords should be the same")
	errUserAlreadyExists  = fmt.Errorf("User with this username already exists")
)

type userRepository interface {
	FindUser(username string) (db.User, error)
	SaveUser(user db.User) error
}

// RegisterHandler struct responsible for handling actions
// made on register html page.
type RegisterHandler struct {
	userRepo  userRepository
	templates *TemplateRepository
}

// NewRegisterHandler returns new RegisterHandler struct.
func NewRegisterHandler(templates *TemplateRepository, userRepo userRepository) *RegisterHandler {
	return &RegisterHandler{
		userRepo:  userRepo,
		templates: templates,
	}
}

// ShowRegisterPage renders register html page.
func (h *RegisterHandler) ShowRegisterPage(w http.ResponseWriter, req *http.Request) {
	RenderTemplate(w, h.templates.Register)
}

// RegisterUser processes user registration form.
func (h *RegisterHandler) RegisterUser(w http.ResponseWriter, req *http.Request) {

	model := NewModel()

	form, err := readRegistrationForm(req)
	if err != nil {
		model.AddError(fmt.Sprintf("Cannot parse form: %v", err))
		RenderTemplateWithModel(w, h.templates.ServerError, model)
		return
	}

	validationErrors := form.validate()

	logger.Infof("Validation errors: %v", validationErrors)

	model["username"] = form.username
	model.AddErrors(validationErrors...)

	if model.HasErrors() {
		RenderTemplateWithModel(w, h.templates.Register, model)
		return
	}

	user, err := h.userRepo.FindUser(form.username)
	if err != nil && err != db.ErrNotFound {
		model.AddError(fmt.Sprintf("Cannot get data about user: %v", err))
		RenderTemplateWithModel(w, h.templates.ServerError, model)
		return
	}

	if !user.Empty() {
		model.AddErrors(errUserAlreadyExists)
		RenderTemplateWithModel(w, h.templates.Register, model)
		return
	}

	passBytes, err := bcrypt.GenerateFromPassword([]byte(form.password1), 1)
	if err != nil {
		model.AddError(fmt.Sprintf("Password encription failed: %v", err))
		RenderTemplateWithModel(w, h.templates.Login, model)
		return
	}

	logger.Infof("Registering user with username: %s", form.username)

	if model.HasErrors() {
		RenderTemplateWithModel(w, h.templates.Register, model)
		return
	}

	user = db.User{Login: form.username, Password: string(passBytes)}
	if err = h.userRepo.SaveUser(user); err != nil {
		model.AddError(fmt.Sprintf("Cannot store user data: %v", err))
		RenderTemplateWithModel(w, h.templates.Login, model)
		return
	}

	http.Redirect(w, req, "/login", http.StatusFound)
}

func readRegistrationForm(req *http.Request) (*registrationForm, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	return &registrationForm{
		username:  req.FormValue("username"),
		password1: req.FormValue("password1"),
		password2: req.FormValue("password2"),
	}, nil
}

type registrationForm struct {
	username  string
	password1 string
	password2 string
}

func (rf *registrationForm) validate() []error {
	errors := make([]error, 0)

	if l := len(rf.username); l < minUsernameLen || l > maxUsernameLen {
		errors = append(errors, errInvalidUsername)
	}

	if l := len(rf.password1); l < minPasswordLen || l > maxPasswordLen {
		errors = append(errors, errInvalidPassword1)
	}

	if l := len(rf.password2); l < minPasswordLen || l > maxPasswordLen {
		errors = append(errors, errInvalidPassword2)
	}

	if rf.password1 != rf.password2 {
		errors = append(errors, errDifferendPasswords)
	}

	return errors
}
