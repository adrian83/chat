package ws

import (
	"chat/logger"
	"sync"
)

const (
	main = "main"
)

// SendError error representing error while sending msg to client.
type SendError struct {
	Client Client
	Err    error
}

// Error implementation of error interface.
func (e SendError) Error() string {
	return e.Error()
}

// Channel is an interface for defining channels.
type Channel interface {
	Name() string
	Clients() map[string]Client
	FindClient(clientID string) (Client, bool)
	SendToEveryone(msg Message) []SendError
	RemoveClient(client Client) []SendError
	AddClient(client Client)
}

// DefaultChannel struct representing chat channel.
type DefaultChannel struct {
	name     string
	lock     sync.RWMutex
	clients  map[string]Client
	channels Channels
}

// Empty returns true if there is no clients in this channel, false otherwise
func (ch *DefaultChannel) Empty() bool {
	ch.lock.RLock()
	result := 0 == len(ch.clients)
	ch.lock.RUnlock()
	return result
}

// FindClient returns client with given id if it exist in this channel.
func (ch *DefaultChannel) FindClient(clientID string) (Client, bool) {
	ch.lock.RLock()
	c, ok := ch.clients[clientID]
	ch.lock.RUnlock()
	return c, ok
}

// Name returns chanel's name.
func (ch *DefaultChannel) Name() string {
	return ch.name
}

// Clients returns chanel's clients.
func (ch *DefaultChannel) Clients() map[string]Client {
	return ch.clients
}

// SendToEveryone sends message to everyone from that channel.
func (ch *DefaultChannel) SendToEveryone(msg Message) []SendError {
	errors := make([]SendError, 0)
	ch.lock.RLock()
	for _, client := range ch.clients {
		logger.Infof("Channel", "SendToEveryone", "Sending msg to %v from channel '%v'.", client, ch.name)
		client.Send(msg)
	}
	ch.lock.RUnlock()
	return errors
}

// AddClient adds client to channel.
func (ch *DefaultChannel) AddClient(client Client) {
	ch.lock.Lock()
	ch.clients[client.ID()] = client
	ch.lock.Unlock()
}

// RemoveClient removes client from the channel.
func (ch *DefaultChannel) RemoveClient(client Client) []SendError {
	if ch.name == main {
		return make([]SendError, 0)
	}

	ch.lock.Lock()
	delete(ch.clients, client.ID())
	ch.lock.Unlock()

	if ch.Empty() {
		ch.channels.RemoveChannel(ch.Name())

		msg := NewRemoveChannelMessage(ch.name)
		mainChannel := ch.channels.GetMainChannel()
		return mainChannel.SendToEveryone(msg)
	}
	return make([]SendError, 0)
}

// NewChannel functions returns new Channel struct.
func NewChannel(name string, client Client, channels Channels) Channel {
	return &DefaultChannel{
		name:     name,
		clients:  map[string]Client{client.ID(): client},
		channels: channels,
	}
}

// MainChannel struct representing main chat channel.
type MainChannel struct {
	DefaultChannel
}

// RemoveClient removes client from the channel.
func (ch *MainChannel) RemoveClient(client Client) []SendError {
	return make([]SendError, 0)
}

// NewMainChannel functions returns new Channel struct.
func NewMainChannel(channels Channels) Channel {
	return &MainChannel{DefaultChannel{
		name:     main,
		clients:  map[string]Client{},
		channels: channels,
	},
	}
}

// AddClient adds client to channel.
func (ch *MainChannel) AddClient(client Client) {
	ch.DefaultChannel.AddClient(client)
	channelNamesMsg := ChannelsNamesMessage(ch.channels.Names())
	client.Send(channelNamesMsg)
}
