package exchange

import (
	"encoding/json"
)

const (
	MsgUserJoinedRoomMT = "USER_JOINED_ROOM"
	MsgUserLeftRoomMT   = "USER_LEFT_ROOM"
	MsgLogoutMT         = "LOGOUT_USER"
	MsgTextMsgMT        = "TEXT_MSG"
	MsgCreateRoomMT     = "CREATE_ROOM"
	MsgRemoveRoomMT     = "REMOVE_ROOM"
	MsgRoomsNamesMT     = "ROOMS_LIST"
	MsgErrorMsgMT       = "ERROR"

	system = "system"
)

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

// NewCreateRoomMessage returns message which can be used for creating new room.
func NewCreateRoomMessage(roomName string) *Message {
	return &Message{
		MsgType:    MsgCreateRoomMT,
		Room:       roomName,
		SenderID:   system,
		SenderName: system,
	}
}

// NewRemoveRoomMessage returns message which can be used for removing room.
func NewRemoveRoomMessage(roomName string) *Message {
	return &Message{
		MsgType:    MsgRemoveRoomMT,
		Room:       roomName,
		SenderID:   system,
		SenderName: system,
	}
}

// RoomsNamesMessage returns message which contains names of all rooms.
func RoomsNamesMessage(roomNames []string) *Message {
	return &Message{
		MsgType:    MsgRoomsNamesMT,
		SenderID:   system,
		SenderName: system,
		Rooms:      roomNames,
	}
}

// ErrorMessage returns error message.
func ErrorMessage(content string) *Message {
	return &Message{
		MsgType:    MsgErrorMsgMT,
		SenderID:   system,
		SenderName: system,
		Content:    content,
	}
}

// NewUserJoinedRoomMessage returns  new UserJoinedRoomMessage message.
func NewUserJoinedRoomMessage(room, senderID string) *Message {
	return &Message{
		MsgType:    MsgUserJoinedRoomMT,
		SenderID:   senderID,
		SenderName: senderID,
		Room:       room,
	}
}

// NewUserLeftRoomMessage returns new UserLeftRoomMessage message.
func NewUserLeftRoomMessage(room, senderID string) *Message {
	return &Message{
		MsgType:    MsgUserLeftRoomMT,
		SenderID:   senderID,
		SenderName: senderID,
		Room:       room,
	}
}
