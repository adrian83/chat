package ws

import (
	"fmt"
)

const (
	// AddChannelMsg symbolizes message that creates new channel.
	AddChannelMsg = "ADD_CH"
)

const (
	system = "system"
)

// Message represents ALL messages exchanged in the app. This may not be the
// best idea, but in such small app maybe it won't be catastrophic. We will see.
type Message struct {
	MsgType    string   `json:"msgType"`
	SenderID   string   `json:"senderId"`
	SenderName string   `json:"senderName"`
	Channels   []string `json:"channels"`
	Channel    string   `json:"channel"`
	Content    string   `json:"content"`
}

// String returns string representation of Message struct.
func (m Message) String() string {
	return fmt.Sprintf("Message {	MsgType: %v, SenderID: %v, SenderName: %v, Channels: %v, Channel: %v, Content: %v }",
		m.MsgType, m.SenderID, m.SenderName, m.Channels, m.Channel, m.Content)
}

// NewAddChannelMessage returns message which can be used for creating new channel.
func NewAddChannelMessage(channelName string) Message {
	return Message{
		MsgType:    AddChannelMsg,
		Channel:    channelName,
		SenderID:   system,
		SenderName: system,
	}
}
