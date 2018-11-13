package client

import (
	"testing"
	"time"

	"github.com/adrian83/chat/chat/db"
	"github.com/adrian83/chat/chat/ws"
	"github.com/adrian83/chat/chat/ws/message"

	"github.com/stretchr/testify/assert"
)

var (
	userData = db.User{
		ID:       "abc-def",
		Login:    "John",
		Password: "secret",
	}

	createRoomMessage = message.NewAddRoomMessage("golang")

	rooms = ws.NewRooms()
)

func TestCreateClientWithGivenId(t *testing.T) {
	clientID := "abcdef-ghijkl"
	client := NewClient(clientID, &userData, rooms, nil)

	assert.Equal(t, clientID, client.ID(), "Invalid client id")
}

func TestMessageHasStringRepresentation(t *testing.T) {
	clientID := "abcdef-ghijkl"
	client := NewClient(clientID, &userData, rooms, nil)

	assert.NotEmpty(t, client.String())
}

func TestSendingMessageOnNotStartedClientShouldDoNothing(t *testing.T) {
	connection := &ConnectionStub{
		sent:     make([]interface{}, 0),
		received: make([]interface{}, 0),
	}
	clientID := "abcdef-ghijkl"
	client := NewClient(clientID, &userData, rooms, connection)

	client.Send(createRoomMessage)

	assert.Equal(t, len(connection.sent), 0, "Invalid number of sent messages")
}

func TestSendingMultioleMessagesOnNotStartedClientShouldStopSender(t *testing.T) {
	connection := &ConnectionStub{
		sent:     make([]interface{}, 0),
		received: make([]interface{}, 0),
	}
	clientID := "abcdef-ghijkl"
	client := NewClient(clientID, &userData, rooms, connection)

	finished := make(chan bool, 5)
	go func(client *Client, finished chan bool) {
		for i := 0; i < 60; i++ {
			client.Send(createRoomMessage)
		}
		finished <- true
	}(client, finished)

	successAfterMs(t, 200, finished)
}

type ConnectionStub struct {
	sendErr    error
	receiveErr error
	closeErr   error
	sent       []interface{}
	received   []interface{}
}

func (c *ConnectionStub) Send(msg interface{}) error {
	if c.sendErr != nil {
		return c.sendErr
	}
	c.sent = append(c.sent, msg)
	return nil
}

func (c *ConnectionStub) Receive(msg interface{}) error {
	if c.receiveErr != nil {
		return c.receiveErr
	}
	c.received = append(c.received, msg)
	return nil
}
func (c *ConnectionStub) Close() error {
	return c.closeErr
}

func successAfterMs(t *testing.T, ms int64, finished chan bool) {
	select {
	case <-finished:
		t.Error("Invalid state")
	case <-time.After(time.Duration(ms) * time.Millisecond):
		t.Log("success")
	}
}
