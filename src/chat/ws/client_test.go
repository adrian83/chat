package ws

import (
	"chat/db"
)

var (
	userData = db.User{
		ID:       "abc-def",
		Name:     "John",
		Password: "secret",
	}

	channels = NewChannels()
)

/*
func TestReturnID(t *testing.T) {
	connection := NewMockConnection()
	clientID := "abcdef-ghijkl"
	client := NewClient(clientID, userData, connection, channels)

	if client.ID() != "abcdef-ghijkl" {
		t.Errorf("Incorrect client ID. Should be %s, but is %s", clientID, client.ID())
	}
}
*/
