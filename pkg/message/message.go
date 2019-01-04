package message

import (
	"encoding/json"

	logger "github.com/sirupsen/logrus"
)

const (
	userJoinedRoomMT = "USER_JOINED_ROOM"
	userLeftRoomMT   = "USER_LEFT_ROOM"
	logoutMT         = "LOGOUT_USER"
	textMsgMT        = "TEXT_MSG"
	createRoomMT     = "CREATE_ROOM"
	removeRoomMT     = "REMOVE_ROOM"
	roomsNamesMT     = "ROOMS_LIST"
	errorMsgMT       = "ERROR"

	system = "system"
)

// Sender defines interface for structs that can send messages.
type Sender interface {
	Send(Message)
	ID() string
}

// Rooms defines interface for structs that can manage clients and rooms.
type Rooms interface {
	AddClientToRoom(string, Sender)
	RemoveClientFromRoom(string, Sender)
	SendMessageOnRoom(Message)
	CreateRoom(string, Sender)
}

// Handler defines interface for structs that can do something with given Sender and Rooms.
type Handler interface {
	DoWith(client Sender, rooms Rooms)
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
	bts, _ := json.Marshal(m)
	return string(bts)
}

// DoWith is an implemetation of Handler interface and defines behaviour of message.
// Different types of messages can use differently given sender and rooms.
func (m *Message) DoWith(client Sender, rooms Rooms) {

	switch m.MsgType {
	case userJoinedRoomMT:
		rooms.AddClientToRoom(m.Room, client)
	case textMsgMT:
		rooms.SendMessageOnRoom(*m)
	case createRoomMT:
		rooms.CreateRoom(m.Room, client)
	case logoutMT:
		client.Send(*m)
	case userLeftRoomMT:
		rooms.RemoveClientFromRoom(m.Room, client)
	default:
		logger.Warnf("Message. Unknown message: '%v'", m.MsgType)
	}

}

// NewCreateRoomMessage returns message which can be used for creating new room.
func NewCreateRoomMessage(roomName string) Message {
	return Message{
		MsgType:    createRoomMT,
		Room:       roomName,
		SenderID:   system,
		SenderName: system,
	}
}

// NewRemoveRoomMessage returns message which can be used for removing room.
func NewRemoveRoomMessage(roomName string) Message {
	return Message{
		MsgType:    removeRoomMT,
		Room:       roomName,
		SenderID:   system,
		SenderName: system,
	}
}

// RoomsNamesMessage returns message which contains names of all rooms.
func RoomsNamesMessage(roomNames []string) Message {
	return Message{
		MsgType:    roomsNamesMT,
		SenderID:   system,
		SenderName: system,
		Rooms:      roomNames,
	}
}

// ErrorMessage returns error message.
func ErrorMessage(content string) Message {
	return Message{
		MsgType:    errorMsgMT,
		SenderID:   system,
		SenderName: system,
		Content:    content,
	}
}

// NewUserJoinedRoomMessage returns  new UserJoinedRoomMessage message.
func NewUserJoinedRoomMessage(room, senderID string) Message {
	return Message{
		MsgType:    userJoinedRoomMT,
		SenderID:   senderID,
		SenderName: senderID,
		Room:       room,
	}
}

// NewUserLeftRoomMessage returns new UserLeftRoomMessage message.
func NewUserLeftRoomMessage(room, senderID string) Message {
	return Message{
		MsgType:    userLeftRoomMT,
		SenderID:   senderID,
		SenderName: senderID,
		Room:       room,
	}
}
