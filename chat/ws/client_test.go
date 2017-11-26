package ws

import (
	"testing"

	"github.com/adrian83/go-chat/chat/db"
)

var (
	userData = db.User{
		ID:       "abc-def",
		Name:     "John",
		Password: "secret",
	}

	channels = NewRooms()
)

func TestReturnID(t *testing.T) {
	clientID := "abcdef-ghijkl"
	client := NewClient(clientID, userData, channels, nil)

	if client.ID() != clientID {
		t.Errorf("Incorrect client ID. Should be %s, but is %s", clientID, client.ID())
	}
}
