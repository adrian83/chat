package ws

import (
	"github.com/adrian83/go-chat/chat/logger"
)

const (
	// Main is the name of the main room.
	Main = "main"
)

// Room is an interface for defining rooms.
type Room interface {
	Name() string
	FindClient(clientID string) Client
	SendToEveryone(msg Message)
	RemoveClient(client Client)
	AddClient(client Client)
}

// DefaultRoom struct representing chat room.
type DefaultRoom struct {
	name             string
	clients          map[string]Client
	rooms            Rooms
	removeClientChan chan Client
	addClientChan    chan Client
	existClient      chan clientExist
	messageChan      chan Message
	interrupt        chan bool
}

func (ch *DefaultRoom) start() {
	logger.Infof("DefaultRoom", "start", "Starting room: '%v'", ch.Name())
mainLoop:
	for {
		select {
		case <-ch.interrupt:
			break mainLoop

		case client := <-ch.removeClientChan:
			delete(ch.clients, client.ID())
			if len(ch.clients) == 0 {
				logger.Infof("DefaultRoom", "start", "Room: '%v' is empty. Should be removed.", ch.Name())
				ch.rooms.RemoveRoom(ch.Name())
			}

		case client := <-ch.addClientChan:
			ch.clients[client.ID()] = client

		case msg := <-ch.messageChan:
			for _, client := range ch.clients {
				logger.Infof("DefaultRoom", "SendToEveryone", "Sending msg to %v from room '%v'.", client, ch.name)
				client.Send(msg)
			}

		case c := <-ch.existClient:
			cli, _ := ch.clients[c.clientID]
			c.existChan <- cli
		}
	}
}

type clientExist struct {
	existChan chan Client
	clientID  string
}

// FindClient returns client with given id if it exist in this room.
func (ch *DefaultRoom) FindClient(clientID string) Client {

	clientChan := make(chan Client, 1)

	ch.existClient <- clientExist{
		existChan: clientChan,
		clientID:  clientID,
	}

	return <-clientChan
}

// Name returns chanel's name.
func (ch *DefaultRoom) Name() string {
	return ch.name
}

// SendToEveryone sends message to everyone from that room.
func (ch *DefaultRoom) SendToEveryone(msg Message) {
	ch.messageChan <- msg
}

// AddClient adds client to room.
func (ch *DefaultRoom) AddClient(client Client) {
	ch.addClientChan <- client
}

// RemoveClient removes client from the room.
func (ch *DefaultRoom) RemoveClient(client Client) {
	ch.removeClientChan <- client
}

// NewRoom functions returns new Room struct.
func NewRoom(name string, rooms Rooms) Room {
	room := &DefaultRoom{
		name:             name,
		clients:          map[string]Client{},
		rooms:            rooms,
		removeClientChan: make(chan Client, 5),
		addClientChan:    make(chan Client, 5),
		messageChan:      make(chan Message, 50),
	}
	go room.start()
	return room
}
