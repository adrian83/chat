package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var (
	sessionIDName      = "session_id"
	defSessionDuration = time.Duration(1000) * time.Minute

	errSessionCookieNotFound = fmt.Errorf("cookie with session id not found")
)

func storeSessionCookie(ID string, w http.ResponseWriter) {

	cookie := &http.Cookie{
		Name:   sessionIDName,
		Value:  ID,
		MaxAge: int(defSessionDuration.Seconds()),
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
