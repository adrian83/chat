package message

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	unknownMessage = Message{
		MsgType:    "blabla",
		SenderID:   "abc",
		SenderName: "John",
		Room:       "golang",
	}

	joinRoomMessage = Message{
		MsgType:    userJoinedRoomMT,
		SenderID:   "abc",
		SenderName: "John",
		Room:       "golang",
	}

	leaveRoomMessage = Message{
		MsgType:    userLeftRoomMT,
		SenderID:   "def",
		SenderName: "Jane",
		Room:       "java",
	}

	textMessage = Message{
		MsgType:    textMsgMT,
		SenderID:   "ghi",
		SenderName: "Steve",
		Room:       "haskell",
		Content:    "Haskell is the best",
	}

	createRoomMessage = Message{
		MsgType:    createRoomMT,
		Room:       "dart",
		SenderID:   "jkl",
		SenderName: "Jessica",
	}

	logoutMessage = Message{
		MsgType:    logoutMT,
		SenderID:   "mno",
		SenderName: "Kurt",
	}
)

func TestMessageHasStringRepresentation(t *testing.T) {
	assert.NotEmpty(t, textMessage.String())
}

func TestMessageDoWithMethod(t *testing.T) {

	var testData = []struct {
		message                      Message
		idExecuted                   int
		sendExecuted                 int
		addClientToRoomExecuted      int
		createRoomExecuted           int
		removeClientFromRoomExecuted int
		sendMessageOnRoomExecuted    int
	}{
		{unknownMessage, 0, 0, 0, 0, 0, 0},
		{joinRoomMessage, 0, 0, 1, 0, 0, 0},
		{leaveRoomMessage, 0, 0, 0, 0, 1, 0},
		{textMessage, 0, 0, 0, 0, 0, 1},
		{createRoomMessage, 0, 0, 0, 1, 0, 0},
		{logoutMessage, 0, 1, 0, 0, 0, 0},
	}

	for _, data := range testData {

		sender := &SenderStub{id: "abc"}
		rooms := &RoomsStub{}

		data.message.DoWith(sender, rooms)

		assert.Equal(t, sender.idExecuted, data.idExecuted, "Invalid number of sender.ID() calls")
		assert.Equal(t, sender.sendExecuted, data.sendExecuted, "Invalid number of sender.Send(Message) calls")

		assert.Equal(t, rooms.addClientToRoomExecuted, data.addClientToRoomExecuted, "Invalid number of rooms.AddClientToRoom(string, Sender) calls")
		assert.Equal(t, rooms.createRoomExecuted, data.createRoomExecuted, "Invalid number of rooms.CreateRoom(string, Sender) calls")
		assert.Equal(t, rooms.removeClientFromRoomExecuted, data.removeClientFromRoomExecuted, "Invalid number of rooms.RemoveClientFromRoom(string, Sender) calls")
		assert.Equal(t, rooms.sendMessageOnRoomExecuted, data.sendMessageOnRoomExecuted, "Invalid number of rooms.SendMessageOnRoom(Message) calls")
	}
}

func TestMessagesConstructors(t *testing.T) {

	roomName := "golang"
	roomNames := []string{"java", "haskell"}
	content := "this is content"
	senderID := "abc-def"

	var testData = []struct {
		message            Message
		expectedType       string
		expectedSenderID   string
		expectedSenderName string
		expectedRooms      []string
		expectedRoom       string
		expectedContent    string
	}{
		{NewCreateRoomMessage(roomName), createRoomMT, system, system, nil, roomName, ""},
		{NewRemoveRoomMessage(roomName), removeRoomMT, system, system, nil, roomName, ""},
		{RoomsNamesMessage(roomNames), roomsNamesMT, system, system, roomNames, "", ""},
		{ErrorMessage(content), errorMsgMT, system, system, nil, "", content},
		{NewUserJoinedRoomMessage(roomName, senderID), userJoinedRoomMT, senderID, senderID, nil, roomName, ""},
		{NewUserLeftRoomMessage(roomName, senderID), userLeftRoomMT, senderID, senderID, nil, roomName, ""},
	}

	for _, data := range testData {

		msg := data.message

		assert.Equal(t, data.expectedType, msg.MsgType, fmt.Sprintf("Invalid MsgType. Expected: %v, actual: %v", data.expectedType, msg.MsgType))
		assert.Equal(t, data.expectedSenderID, msg.SenderID, fmt.Sprintf("Invalid SenderID. Expected: %v, actual: %v", data.expectedSenderID, msg.SenderID))
		assert.Equal(t, data.expectedSenderName, msg.SenderName, fmt.Sprintf("Invalid SenderName. Expected: %v, actual: %v", data.expectedSenderName, msg.SenderName))
		assert.Equal(t, data.expectedRooms, msg.Rooms, fmt.Sprintf("Invalid Rooms. Expected: %v, actual: %v", data.expectedRooms, msg.Rooms))
		assert.Equal(t, data.expectedRoom, msg.Room, fmt.Sprintf("Invalid Room. Expected: %v, actual: %v", data.expectedRoom, msg.Room))
		assert.Equal(t, data.expectedContent, msg.Content, fmt.Sprintf("Invalid Content. Expected: %v, actual: %v", data.expectedContent, msg.Content))
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
