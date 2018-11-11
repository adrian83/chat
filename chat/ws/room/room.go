package room

import (
	"fmt"

	"github.com/adrian83/chat/chat/logger"
	"github.com/adrian83/chat/chat/ws/message"
)

const (
	// main is the name of the main room.
	main = "main"
)

// MainRoomName returns name of the main room.
func MainRoomName() string {
	return main
}

// Rooms is an interface for all structs which can manage Room structs.
type Rooms interface {
	RemoveRoom(string)
}

// NewRoom functions returns new Room struct.
func NewRoom(name string, rooms Rooms) *Room {
	room := &Room{
		name:             name,
		clients:          map[string]message.Sender{},
		rooms:            rooms,
		removeClientChan: make(chan string, 5),
		addClientChan:    make(chan message.Sender, 5),
		incomingMessages: make(chan message.Message, 50),
	}
	return room
}

// NewMainRoom returns new unremovable Room struct with name 'main'.
func NewMainRoom(rooms Rooms) *Room {
	return NewRoom(main, rooms)
}

// Room represents chat room.
type Room struct {
	name             string
	clients          map[string]message.Sender
	rooms            Rooms
	removeClientChan chan string
	addClientChan    chan message.Sender
	clientExists     chan clientExist
	incomingMessages chan message.Message
	interrupt        chan bool
}

// FindClient returns client with given id if it exist in this room.
func (ch *Room) FindClient(clientID string) (message.Sender, error) {

	clientChan := make(chan message.Sender, 1)

	ch.clientExists <- clientExist{
		existChan: clientChan,
		clientID:  clientID,
	}

	sender := <-clientChan
	if sender == nil {
		return nil, fmt.Errorf("client with id %v cannot be found", clientID)
	}
	return sender, nil
}

// Main returns true if this room is a main room.
func (ch *Room) Main() bool {
	return ch.name == main
}

// Name returns room's name.
func (ch *Room) Name() string {
	return ch.name
}

// SendToEveryone sends message to everyone in this room.
func (ch *Room) SendToEveryone(msg message.Message) {
	ch.incomingMessages <- msg
}

// AddClient adds client to this room.
func (ch *Room) AddClient(client message.Sender) {
	ch.addClientChan <- client
}

// RemoveClient removes client from this room.
func (ch *Room) RemoveClient(clientID string) {
	ch.removeClientChan <- clientID
}

// Start starts room. After invoking this method room can process sent messages.
func (ch *Room) Start() {
	go func() {
		for {
			select {
			case <-ch.interrupt:
				return

			case clientID := <-ch.removeClientChan:
				delete(ch.clients, clientID)
				if len(ch.clients) == 0 {
					logger.Infof("DefaultRoom", "start", "Room: '%v' is empty. Should be removed.", ch.Name())
					ch.rooms.RemoveRoom(ch.Name())
				}

			case client := <-ch.addClientChan:
				ch.clients[client.ID()] = client

			case msg := <-ch.incomingMessages:
				for _, client := range ch.clients {
					logger.Infof("DefaultRoom", "SendToEveryone", "Sending msg to %v from room '%v'.", client, ch.name)
					client.Send(msg)
				}

			case c := <-ch.clientExists:
				c.existChan <- ch.clients[c.clientID]
			}
		}
	}()
}

type clientExist struct {
	existChan chan message.Sender
	clientID  string
}
