package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adrian83/chat/pkg/config"
	"github.com/adrian83/chat/pkg/user"
)

type userRepositoryMock struct {
	findErr      error
	saveErr      error
	user         *user.User
	saveExecuted bool
	findExecuted bool
}

func (ur *userRepositoryMock) FindUser(username string) (*user.User, error) {
	return ur.user, ur.findErr
}

func (ur *userRepositoryMock) SaveUser(user user.User) error {
	return ur.saveErr
}

func TestFormValidation(t *testing.T) {
	config := config.StaticsConfig{Path: "../../static"}
	templatesRepository := NewTemplateRepository(config)

	userRepository := &userRepositoryMock{user: &user.User{ID: "abc", Login: "John", Password: "dyuwqgdjhqg"}}

	handler := NewRegisterHandler(templatesRepository, userRepository)

	var testData = []struct {
		form        *registrationForm
		presentErrs []error
		missingErrs []error
	}{
		{
			form:        &registrationForm{username: "John", password1: "secret", password2: "secret"},
			presentErrs: nil,
			missingErrs: []error{ErrInvalidUsername, ErrInvalidPassword1, ErrInvalidPassword2, ErrDifferendPasswords},
		},
		{
			form:        &registrationForm{username: "", password1: "secret", password2: "secret"},
			presentErrs: []error{ErrInvalidUsername},
			missingErrs: []error{ErrInvalidPassword1, ErrInvalidPassword2, ErrDifferendPasswords},
		},
		{
			form:        &registrationForm{username: "John", password1: "", password2: "secret"},
			presentErrs: []error{ErrInvalidPassword1, ErrDifferendPasswords},
			missingErrs: []error{ErrInvalidUsername, ErrInvalidPassword2},
		},
		{
			form:        &registrationForm{username: "John", password1: "secret", password2: ""},
			presentErrs: []error{ErrInvalidPassword2, ErrDifferendPasswords},
			missingErrs: []error{ErrInvalidUsername, ErrInvalidPassword1},
		},
		{
			form:        &registrationForm{username: "John", password1: "secret", password2: "other"},
			presentErrs: []error{ErrDifferendPasswords},
			missingErrs: []error{ErrInvalidUsername, ErrInvalidPassword1, ErrInvalidPassword2},
		},
		{
			form:        &registrationForm{username: "", password1: "", password2: ""},
			presentErrs: []error{ErrInvalidUsername, ErrInvalidPassword1, ErrInvalidPassword2},
			missingErrs: []error{ErrDifferendPasswords},
		},
		{
			form:        &registrationForm{username: "", password1: "b", password2: "a"},
			presentErrs: []error{ErrInvalidUsername, ErrInvalidPassword1, ErrInvalidPassword2, ErrDifferendPasswords},
			missingErrs: nil,
		},
	}

	for _, data := range testData {

		request := httptest.NewRequest(http.MethodPost, "/test", nil)
		response := httptest.NewRecorder()

		request.PostForm = map[string][]string{
			"username":  []string{data.form.username},
			"password1": []string{data.form.password1},
			"password2": []string{data.form.password2},
		}

		handler.RegisterUser(response, request)

		respString := responseToString(response.Result())

		for _, err := range data.presentErrs {
			assert.Contains(t, respString, err.Error())
		}

		for _, err := range data.missingErrs {
			assert.NotContains(t, respString, err.Error())
		}
	}
}

func TestRegistrationShouldFailIfUserRepositoryCannotCheckIfUsernameIsUnique(t *testing.T) {
	config := config.StaticsConfig{Path: "../../static"}
	templatesRepository := NewTemplateRepository(config)
	err := fmt.Errorf("error while getting user data from db")

	userRepository := &userRepositoryMock{findErr: err}

	handler := NewRegisterHandler(templatesRepository, userRepository)

	request := httptest.NewRequest(http.MethodPost, "/test", nil)
	request.PostForm = map[string][]string{
		"username":  []string{"John"},
		"password1": []string{"secret"},
		"password2": []string{"secret"},
	}

	response := httptest.NewRecorder()

	handler.RegisterUser(response, request)

	respString := responseToString(response.Result())

	assert.Contains(t, respString, err.Error())
}

func responseToString(response *http.Response) string {
	var buf bytes.Buffer
	buf.ReadFrom(response.Body)
	return buf.String()
}
