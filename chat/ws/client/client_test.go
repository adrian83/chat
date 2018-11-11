package client

import (
	"testing"

	"github.com/adrian83/chat/chat/db"
	"github.com/adrian83/chat/chat/ws"

	"github.com/stretchr/testify/assert"
)

var (
	userData = db.User{
		ID:       "abc-def",
		Login:    "John",
		Password: "secret",
	}

	rooms = ws.NewRooms()
)

func TestReturnID(t *testing.T) {
	clientID := "abcdef-ghijkl"
	client := NewClient(clientID, &userData, rooms, nil)

	assert.Equal(t, clientID, client.ID(), "Invalid client id")
}
