package ws

import (
	"fmt"
)

const (
	// UserJoinedRoomMT is a marker for messages send when user is joining a room.
	UserJoinedRoomMT = "USER_JOINED_ROOM"
	// UserLeftRoomMT is a marker for messages send when user is leaving a room.
	UserLeftRoomMT = "USER_LEFT_ROOM"
	// LogoutMT is a marker for logout messages.
	LogoutMT = "LOGOUT_USER"
	// TextMsgMT is a marker for text messages.
	TextMsgMT = "TEXT_MSG"
	// CreateRoomMT is a marker for messages creating new room.
	CreateRoomMT = "CREATE_ROOM"
	// RemoveRoomMT is a marker for messages creating new room.
	RemoveRoomMT = "REMOVE_ROOM"
	// RoomsNamesMT is a marker for messages with list of all existing rooms.
	RoomsNamesMT = "ROOMS_LIST"
	// ErrorMsgMT is a marker for error messages.
	ErrorMsgMT = "ERROR"
)

const (
	system = "system"
)

// Doer does something.
type Doer interface {
	Do(client Client, rooms Rooms)
}

// Message represents ALL messages exchanged in the app. This may not be the
// best idea, but in such small app maybe it won't be catastrophic. We will see.
type Message struct {
	MsgType    string   `json:"msgType"`
	SenderID   string   `json:"senderId"`
	SenderName string   `json:"senderName"`
	Rooms      []string `json:"rooms"`
	Room       string   `json:"room"`
	Content    string   `json:"content"`
}

// String returns string representation of Message struct.
func (m *Message) String() string {
	return fmt.Sprintf("Message {	MsgType: %v, SenderID: %v, SenderName: %v, Rooms: %v, Room: %v, Content: %v }",
		m.MsgType, m.SenderID, m.SenderName, m.Rooms, m.Room, m.Content)
}

// Do handles message.
func (m *Message) Do(client Client, rooms Rooms) {
	var doer Doer
	switch m.MsgType {
	case TextMsgMT:
		doer = &TextMessage{Message: m}
	case CreateRoomMT:
		doer = &NewRoomMessage{Message: m}
	case LogoutMT:
		doer = &LogoutMessage{Message: m}
	case UserJoinedRoomMT:
		doer = &JoinRoomMessage{Message: m}
	case UserLeftRoomMT:
		doer = &LeaveRoomMessage{Message: m}
	default:
		doer = &UnknownMessage{Message: m}
	}
	doer.Do(client, rooms)

}

type UnknownMessage struct {
	*Message
}

func (m *UnknownMessage) Do(client Client, rooms Rooms) {
	//rooms.SendMessageOnRoom(*m.Message)
}

type JoinRoomMessage struct {
	*Message
}

func (m *JoinRoomMessage) Do(client Client, rooms Rooms) {
	rooms.AddClientToRoom(m.Room, client)
}

type LeaveRoomMessage struct {
	*Message
}

func (m *LeaveRoomMessage) Do(client Client, rooms Rooms) {
	rooms.RemoveClientFromRoom(m.Room, client)
}

type TextMessage struct {
	*Message
}

func (m *TextMessage) Do(client Client, rooms Rooms) {
	rooms.SendMessageOnRoom(*m.Message)
}

type NewRoomMessage struct {
	*Message
}

func (m *NewRoomMessage) Do(client Client, rooms Rooms) {
	rooms.CreateRoom(m.Room, client)
}

type LogoutMessage struct {
	*Message
}

func (m *LogoutMessage) Do(client Client, rooms Rooms) {
	client.Send(*m.Message)
}

// NewAddRoomMessage returns message which can be used for creating new room.
func NewAddRoomMessage(roomName string) Message {
	return Message{
		MsgType:    CreateRoomMT,
		Room:       roomName,
		SenderID:   system,
		SenderName: system,
	}
}

// NewRemoveRoomMessage returns message which can be used for removing room.
func NewRemoveRoomMessage(roomName string) Message {
	return Message{
		MsgType:    RemoveRoomMT,
		Room:       roomName,
		SenderID:   system,
		SenderName: system,
	}
}

// RoomsNamesMessage returns message which contains names of all rooms.
func RoomsNamesMessage(roomNames []string) Message {
	return Message{
		MsgType:    RoomsNamesMT,
		SenderID:   system,
		SenderName: system,
		Rooms:      roomNames,
	}
}

// ErrorMessage returns error message.
func ErrorMessage(content string) Message {
	return Message{
		MsgType:    ErrorMsgMT,
		SenderID:   system,
		SenderName: system,
		Content:    content,
	}
}

// NewUserJoinedRoomMessage returns  new UserJoinedRoomMessage message.
func NewUserJoinedRoomMessage(room, senderID string) Message {
	return Message{
		MsgType:    UserJoinedRoomMT,
		SenderID:   senderID,
		SenderName: senderID,
		Room:       room,
	}
}

// NewUserLeftRoomMessage returns new UserLeftRoomMessage message.
func NewUserLeftRoomMessage(room, senderID string) Message {
	return Message{
		MsgType:    UserLeftRoomMT,
		SenderID:   senderID,
		SenderName: senderID,
		Room:       room,
	}
}
