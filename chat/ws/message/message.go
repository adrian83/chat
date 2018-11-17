package message

import (
	"encoding/json"
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
	var handler Handler
	switch m.MsgType {
	case textMsgMT:
		handler = &TextMessage{Message: m}
	case createRoomMT:
		handler = &CreateRoomMessage{Message: m}
	case logoutMT:
		handler = &LogoutMessage{Message: m}
	case userJoinedRoomMT:
		handler = &JoinRoomMessage{Message: m}
	case userLeftRoomMT:
		handler = &LeaveRoomMessage{Message: m}
	default:
		handler = &UnknownMessage{Message: m}
	}
	handler.DoWith(client, rooms)

}

// UnknownMessage represents message with unknown MsgType.
type UnknownMessage struct {
	*Message
}

// DoWith defines behaviour of message. Since this message is unknown to system it will be only logged.
func (m *UnknownMessage) DoWith(client Sender, rooms Rooms) {
	// TODO add log here
}

// JoinRoomMessage represents messages which should result in adding client to room (client and room data are inside of the message).
type JoinRoomMessage struct {
	*Message
}

// DoWith defines behaviour of message. Adds client to room.
func (m *JoinRoomMessage) DoWith(client Sender, rooms Rooms) {
	rooms.AddClientToRoom(m.Room, client)
}

// LeaveRoomMessage represents messages which should result in removing client from room (client and room data are inside of the message).
type LeaveRoomMessage struct {
	*Message
}

// DoWith defines behaviour of message. Removes client from room.
func (m *LeaveRoomMessage) DoWith(client Sender, rooms Rooms) {
	rooms.RemoveClientFromRoom(m.Room, client)
}

// TextMessage represents messages which should result in sending text to all clients from given room (room name is inside of the message).
type TextMessage struct {
	*Message
}

// DoWith defines behaviour of message. Sends text to all clients from room with name in message.
func (m *TextMessage) DoWith(client Sender, rooms Rooms) {
	rooms.SendMessageOnRoom(*m.Message)
}

// CreateRoomMessage represents messages which should result in creating new room (room name is inside of the message).
type CreateRoomMessage struct {
	*Message
}

// DoWith defines behaviour of message. Creates new room.
func (m *CreateRoomMessage) DoWith(client Sender, rooms Rooms) {
	rooms.CreateRoom(m.Room, client)
}

// LogoutMessage represents messages which should result in removing user from all rooms.
type LogoutMessage struct {
	*Message
}

// DoWith defines behaviour of message. User loggout.
func (m *LogoutMessage) DoWith(client Sender, rooms Rooms) {
	client.Send(*m.Message)
}

// --------------------------------

// NewCreateRoomMessage returns message which can be used for creating new room.
func NewCreateRoomMessage(roomName string) *CreateRoomMessage {
	return &CreateRoomMessage{&Message{
		MsgType:    createRoomMT,
		Room:       roomName,
		SenderID:   system,
		SenderName: system,
	}}
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
