package handler

import (
	"net/http"
	"time"
)

var (
	sessionIDName      = "session_id"
	defSessionDuration = time.Duration(1000) * time.Minute
)

func storeSessionCookie(ID string, w http.ResponseWriter) {

	cookie := &http.Cookie{
		Name:   sessionIDName,
		Value:  ID,
		MaxAge: int(defSessionDuration.Seconds()),
	}

	http.SetCookie(w, cookie)
}
