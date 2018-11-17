package client

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/adrian83/chat/chat/ws/message"
	"github.com/adrian83/chat/chat/ws/rooms"

	"github.com/stretchr/testify/assert"
)

var (
	createRoomMessage = message.NewCreateRoomMessage("golang")

	chatRooms = rooms.NewRooms()
)

func TestCreateClientWithGivenId(t *testing.T) {
	clientID := "abcdef-ghijkl"
	user := &UserStub{name: "John"}

	client := NewClient(clientID, user, chatRooms, nil)

	assert.Equal(t, clientID, client.ID(), "Invalid client id")
}

func TestMessageHasStringRepresentation(t *testing.T) {
	clientID := "abcdef-ghijkl"
	user := &UserStub{name: "John"}

	client := NewClient(clientID, user, chatRooms, nil)

	assert.NotEmpty(t, client.String())
}

func TestSendingMessageOnNotStartedClientShouldDoNothing(t *testing.T) {
	connection := &ConnectionStub{t: t}
	clientID := "abcdef-ghijkl"
	user := &UserStub{name: "John"}

	client := NewClient(clientID, user, chatRooms, connection)

	client.Send(createRoomMessage)

	assert.Equal(t, len(connection.sent), 0, "Invalid number of sent messages")
	assert.Equal(t, 1, user.executed, "User should be asked for name once (during client creation)")
}

func TestSendingMultipleMessagesOnNotStartedClientShouldStopSender(t *testing.T) {
	connection := &ConnectionStub{t: t}
	clientID := "abcdef-ghijkl"
	user := &UserStub{name: "John"}

	client := NewClient(clientID, user, chatRooms, connection)

	finished := make(chan bool, 5)
	go func(client *Client, finished chan bool) {
		for i := 0; i < 60; i++ {
			client.Send(createRoomMessage)
		}
		finished <- true
	}(client, finished)

	successAfterMs(t, 200, finished)
}

func TestStartingClientShouldSuspendExecution(t *testing.T) {
	connection := &ConnectionStub{t: t}
	clientID := "abcdef-ghijkl"
	user := &UserStub{name: "John"}

	client := NewClient(clientID, user, chatRooms, connection)

	finished := make(chan bool, 5)
	go func(client *Client, finished chan bool) {
		client.Start()
		finished <- true
	}(client, finished)

	successAfterMs(t, 200, finished)
}

func TestErrorWhileReceivingMessageShouldStopTheClient(t *testing.T) {
	receivedMsg := msgOrErr{
		msg: createRoomMessage,
	}

	receivedErr := msgOrErr{
		err: fmt.Errorf("test error"),
	}

	connection := &ConnectionStub{
		t:        t,
		received: []msgOrErr{receivedMsg, receivedErr},
	}
	clientID := "abcdef-ghijkl"
	user := &UserStub{name: "John"}

	client := NewClient(clientID, user, chatRooms, connection)

	finished := make(chan bool, 5)
	go failAfterMs(t, 200, finished)

	client.Start()
	finished <- true
}

func TestClientShouldBeStoppedAfterCallingStopMethod(t *testing.T) {

	connection := &ConnectionStub{t: t}
	clientID := "abcdef-ghijkl"
	user := &UserStub{name: "John"}

	client := NewClient(clientID, user, chatRooms, connection)

	finished := make(chan bool, 5)
	go failAfterMs(t, 200, finished)

	client.stop()
	client.Start()
	finished <- true
}

func TestHandleErrorWhileClosingConnection(t *testing.T) {
	connection := &ConnectionStub{
		t:        t,
		closeErr: fmt.Errorf("test"),
	}
	clientID := "abcdef-ghijkl"
	user := &UserStub{name: "John"}

	client := NewClient(clientID, user, chatRooms, connection)

	finished := make(chan bool, 5)

	client.stop()
	client.Start()
	finished <- true
}

func TestErrorWhileSendingMessageShouldStopTheClient(t *testing.T) {
	connection := &ConnectionStub{
		t: t,
		sendErr: &sendErr{
			no:  1,
			err: fmt.Errorf("test"),
		},
	}
	clientID := "abcdef-ghijkl"
	user := &UserStub{name: "John"}

	client := NewClient(clientID, user, chatRooms, connection)

	client.Send(createRoomMessage)
	client.Send(createRoomMessage)

	finished := make(chan bool, 5)

	go func(client *Client, finished chan bool) {
		time.Sleep(time.Duration(100) * time.Millisecond)
		client.stop()
		finished <- true
	}(client, finished)

	go failAfterMs(t, 200, finished)

	client.Start()
}

type UserStub struct {
	name     string
	executed int
}

func (u *UserStub) Name() string {
	u.executed++
	return u.name
}

type msgOrErr struct {
	msg message.Message
	err error
}

type sendErr struct {
	no  int
	err error
}

type ConnectionStub struct {
	t        *testing.T
	sendErr  *sendErr
	closeErr error
	sent     []interface{}
	received []msgOrErr
	wait     chan interface{}
}

func (c *ConnectionStub) Send(msg interface{}) error {
	c.sent = append(c.sent, msg)

	if c.sendErr != nil {
		if c.sendErr.no == 0 && c.sendErr.err != nil {
			return c.sendErr.err
		}
		c.sendErr.no--
	}

	return nil
}

func (c *ConnectionStub) Receive(msg interface{}) error {
	if len(c.received) == 0 {
		<-c.wait
	}

	msgOrErr := c.received[0]
	c.received = c.received[1:]

	if msgOrErr.err != nil {
		return msgOrErr.err
	}

	bts, err := json.Marshal(msgOrErr.msg)
	if err != nil {
		c.t.Errorf("unexpected error while marshaling message to json")
	}

	if err = json.Unmarshal(bts, msg); err != nil {
		c.t.Errorf("unexpected error while unmarshaling message from json")
	}

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

func failAfterMs(t *testing.T, ms int64, finished chan bool) {
	select {
	case <-finished:
		t.Log("success")
	case <-time.After(time.Duration(ms) * time.Millisecond):
		t.Error("Invalid state")
	}
}
