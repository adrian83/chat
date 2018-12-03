package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
func ReadSessionIDFromCookie(req *http.Request) (string, error) {
	sessionCookie, err := req.Cookie(sessionIDName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", errSessionCookieNotFound
		}
		return "", errors.Wrap(err, "error while getting sessionID from request")
	}
	return sessionCookie.Value, nil
}

*/

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
