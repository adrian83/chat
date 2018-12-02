package session

import (
	"github.com/adrian83/chat/chat/db"

	redisSession "github.com/adrian83/go-redis-session"
	"github.com/google/uuid"

	"net/http"
	"time"
)

const (
	sessionIDName      = "session_id"
	defSessionDuration = time.Duration(1000) * time.Minute
)

// FindSessionID returns session ID from cookie or empty string if such
// cookie doesn't exist, error if something bad has happened.
func FindSessionID(req *http.Request) string {
	c, err := req.Cookie(sessionIDName)
	if err != nil {
		return ""
	}
	return c.Value
}

// Session struct represents simplified session mechanism.
type Session struct {
	sessionStore *redisSession.Store
}

// New returns pointer to new Session struct.
func New(sessionStore *redisSession.Store) *Session {
	return &Session{sessionStore: sessionStore}
}

// FindUserData returns user data with given sessionID from session or empty
// struct if session doesn't contain such data, error if something
// bad has happened.
func (s *Session) FindUserData(sessionID string) (db.User, error) {
	session, err := s.sessionStore.Find(sessionID)
	if err != nil {
		return db.User{}, err
	}

	user := new(db.User)
	if err = session.Get("user", user); err != nil {
		return db.User{}, nil
	}

	return *user, nil
}

// StoreUserData saves user data to session and returns session ID, error
// if something bad has happened.
func (s *Session) StoreUserData(w http.ResponseWriter, user db.User) (string, error) {
	sessionID := uuid.New().String()
	session, err := s.sessionStore.Create(sessionID, defSessionDuration)
	if err != nil {
		return "", err
	}
	err = session.Add("user", user)
	if err != nil {
		return "", err
	}

	if err := s.sessionStore.Save(session); err != nil {
		return "", err
	}

	cookie := &http.Cookie{
		Name:   sessionIDName,
		Value:  sessionID,
		MaxAge: int(defSessionDuration.Seconds()),
	}

	http.SetCookie(w, cookie)

	return sessionID, nil
}

// Remove removes the sessionID from cookie and session from database.
func (s *Session) Remove(w http.ResponseWriter, req *http.Request) error {
	c, err := req.Cookie(sessionIDName)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:   sessionIDName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	return s.sessionStore.Delete(c.Value)
}
