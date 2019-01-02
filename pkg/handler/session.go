package handler

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// SessionValidFor represents time for how long session is stored
// in redis and for how long cookie with session id exists.
const SessionValidFor = 3600 // seconds

var (
	sessionIDName = "session_id"

	errSessionCookieNotFound = fmt.Errorf("cookie with session id not found")
)

// StoreSessionCookie stores session cookie with given session id.
func StoreSessionCookie(ID string, w http.ResponseWriter) {

	cookie := &http.Cookie{
		Name:   sessionIDName,
		Value:  ID,
		MaxAge: SessionValidFor,
	}

	http.SetCookie(w, cookie)
}

// RemoveSessionCookie removes cookie wirt session id.
func RemoveSessionCookie(w http.ResponseWriter) {

	cookie := &http.Cookie{
		Name:   sessionIDName,
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

// ReadSessionIDFromCookie returns session id read from cookie or error if id cannot be found.
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
