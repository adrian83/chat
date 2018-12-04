package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadingSessionIdFromCookie(t *testing.T) {
	cookie := &http.Cookie{
		Name:   sessionIDName,
		Value:  "abc",
		MaxAge: 30,
	}
	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	request.AddCookie(cookie)

	val, err := ReadSessionIDFromCookie(request)

	assert.Equal(t, cookie.Value, val, "Invalid cookie value")
	assert.NoError(t, err, "Unexpected error while reading session cookie")
}
func TestReadingUnexistingSessionIdFromCookie(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/test", nil)

	val, err := ReadSessionIDFromCookie(request)

	assert.Equal(t, "", val, "Invalid cookie value")
	assert.EqualError(t, err, errSessionCookieNotFound.Error(), "Unexpected error type")
}
