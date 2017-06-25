package ws

import (
	"fmt"
)

const (
	// AddChannelMsg symbolizes message that creates new channel.
	AddChannelMsg = "ADD_CH"
	// RemoveChannelMsg symbolizes message that removes channel.
	RemoveChannelMsg = "REM_CH"
	// ChannelsNamesMsg symbolizes message that contains names of all channel.
	ChannelsNamesMsg = "CHAN_LIST_MSG"
	// ErrorMsg symbolizes error message.
	ErrorMsg = "ERROR"
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

// NewRemoveChannelMessage returns message which can be used for removing channel.
func NewRemoveChannelMessage(channelName string) Message {
	return Message{
		MsgType:    RemoveChannelMsg,
		Channel:    channelName,
		SenderID:   system,
		SenderName: system,
	}
}

// ChannelsNamesMessage returns message which contains names of all channels.
func ChannelsNamesMessage(channelNames []string) Message {
	return Message{
		MsgType:    ChannelsNamesMsg,
		SenderID:   system,
		SenderName: system,
		Channels:   channelNames,
	}
}

// ErrorMessage returns error message.
func ErrorMessage(content string) Message {
	return Message{
		MsgType:    ErrorMsg,
		SenderID:   system,
		SenderName: system,
		Content:    content,
	}
}
