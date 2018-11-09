package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	unknownMessage = UnknownMessage{&Message{
		MsgType:    "blabla",
		SenderID:   "abc",
		SenderName: "John",
		Room:       "golang",
	}}

	joinRoomMessage = JoinRoomMessage{&Message{
		MsgType:    userJoinedRoomMT,
		SenderID:   "abc",
		SenderName: "John",
		Room:       "golang",
	}}

	leaveRoomMessage = LeaveRoomMessage{&Message{
		MsgType:    userLeftRoomMT,
		SenderID:   "def",
		SenderName: "Jane",
		Room:       "java",
	}}
)

func TestUnknownMessageImplementsHandler(t *testing.T) {
	var _ Handler = &UnknownMessage{}
}

func TestJoinRoomMessageImplementsHandler(t *testing.T) {
	var _ Handler = &JoinRoomMessage{}
}

func TestLeaveRoomMessageImplementsHandler(t *testing.T) {
	var _ Handler = &LeaveRoomMessage{}
}

func TestHandlerShouldInteractWithGivenRoomsAndSender(t *testing.T) {

	var testData = []struct {
		handler                      Handler
		idExecuted                   int
		sendExecuted                 int
		addClientToRoomExecuted      int
		createRoomExecuted           int
		removeClientFromRoomExecuted int
		sendMessageOnRoomExecuted    int
	}{
		{&unknownMessage, 0, 0, 0, 0, 0, 0},
		{&joinRoomMessage, 0, 0, 1, 0, 0, 0},
		{&leaveRoomMessage, 0, 0, 0, 0, 1, 0},
	}

	for _, data := range testData {

		sender := &SenderStub{id: "abc"}
		rooms := &RoomsStub{}

		data.handler.DoWith(sender, rooms)

		assert.Equal(t, sender.idExecuted, data.idExecuted, "Invalid number of sender.ID() calls")
		assert.Equal(t, sender.sendExecuted, data.sendExecuted, "Invalid number of sender.Send(Message) calls")

		assert.Equal(t, rooms.addClientToRoomExecuted, data.addClientToRoomExecuted, "Invalid number of rooms.AddClientToRoom(string, Sender) calls")
		assert.Equal(t, rooms.createRoomExecuted, data.createRoomExecuted, "Invalid number of rooms.RemoveClientFromRoom(string, Sender) calls")
		assert.Equal(t, rooms.removeClientFromRoomExecuted, data.removeClientFromRoomExecuted, "Invalid number of rooms.SendMessageOnRoom(Message) calls")
		assert.Equal(t, rooms.sendMessageOnRoomExecuted, data.sendMessageOnRoomExecuted, "Invalid number of rooms.CreateRoom(string, Sender) calls")
	}
}

type SenderStub struct {
	id           string
	sendExecuted int
	idExecuted   int
}

func (s *SenderStub) Send(msg Message) {
	s.sendExecuted++
}

func (s *SenderStub) ID() string {
	s.idExecuted++
	return s.id
}

type RoomsStub struct {
	addClientToRoomExecuted      int
	removeClientFromRoomExecuted int
	sendMessageOnRoomExecuted    int
	createRoomExecuted           int
}

func (r *RoomsStub) AddClientToRoom(string, Sender) {
	r.addClientToRoomExecuted++
}

func (r *RoomsStub) RemoveClientFromRoom(string, Sender) {
	r.removeClientFromRoomExecuted++
}

func (r *RoomsStub) SendMessageOnRoom(Message) {
	r.sendMessageOnRoomExecuted++
}

func (r *RoomsStub) CreateRoom(string, Sender) {
	r.createRoomExecuted++
}
