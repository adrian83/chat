package ws

import (
	"chat/logger"
	"sync"
)

const (
	// Main is the name of the main room.
	Main = "main"
)

// Channel is an interface for defining channels.
type Channel interface {
	Name() string
	FindClient(clientID string) (Client, bool)
	SendToEveryone(msg Message)
	RemoveClient(client Client)
	AddClient(client Client)
}

// DefaultChannel struct representing chat channel.
type DefaultChannel struct {
	name             string
	lock             sync.RWMutex
	clients          map[string]Client
	channels         Channels
	removeClientChan chan Client
	messageChan      chan Message
	interrupt        chan bool
}

func (ch *DefaultChannel) start() {
	logger.Infof("DefaultChannel", "start", "Starting room: '%v'", ch.Name())
mainLoop:
	for {
		select {
		case <-ch.interrupt:
			break mainLoop

		case client := <-ch.removeClientChan:
			delete(ch.clients, client.ID())
			if len(ch.clients) == 0 {
				logger.Infof("DefaultChannel", "start", "Room: '%v' is empty. Should be removed.", ch.Name())
				ch.channels.RemoveChannel(ch.Name())
			}
		case msg := <-ch.messageChan:
			for _, client := range ch.clients {
				logger.Infof("Channel", "SendToEveryone", "Sending msg to %v from channel '%v'.", client, ch.name)
				client.Send(msg)
			}

		}
	}

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

// SendToEveryone sends message to everyone from that channel.
func (ch *DefaultChannel) SendToEveryone(msg Message) {
	ch.messageChan <- msg
}

// AddClient adds client to channel.
func (ch *DefaultChannel) AddClient(client Client) {
	ch.lock.Lock()
	ch.clients[client.ID()] = client
	ch.lock.Unlock()
}

// RemoveClient removes client from the channel.
func (ch *DefaultChannel) RemoveClient(client Client) {
	ch.removeClientChan <- client
}

// NewChannel functions returns new Channel struct.
func NewChannel(name string, client Client, channels Channels) Channel {
	room := &DefaultChannel{
		name:             name,
		clients:          map[string]Client{client.ID(): client},
		channels:         channels,
		removeClientChan: make(chan Client, 5),
		messageChan:      make(chan Message, 50),
	}

	go room.start()

	return room
}

// MainChannel struct representing main chat channel.
type MainChannel struct {
	DefaultChannel
}

// NewMainChannel functions returns new Channel struct.
func NewMainChannel(channels Channels) Channel {
	room := &MainChannel{DefaultChannel{
		name:             Main,
		clients:          map[string]Client{},
		channels:         channels,
		removeClientChan: make(chan Client, 5),
		messageChan:      make(chan Message, 50),
	},
	}

	go room.start()

	return room
}

// AddClient adds client to channel.
func (ch *MainChannel) AddClient(client Client) {
	ch.DefaultChannel.AddClient(client)
}
