package room

import (
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

func (ch *Room) start() {
	logger.Infof("DefaultRoom", "start", "Starting room: '%v'", ch.Name())
mainLoop:
	for {
		select {
		case <-ch.interrupt:
			break mainLoop

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
			cli, _ := ch.clients[c.clientID]
			c.existChan <- cli
		}
	}
}

type clientExist struct {
	existChan chan message.Sender
	clientID  string
}

// FindClient returns client with given id if it exist in this room.
func (ch *Room) FindClient(clientID string) message.Sender {

	clientChan := make(chan message.Sender, 1)

	ch.clientExists <- clientExist{
		existChan: clientChan,
		clientID:  clientID,
	}

	return <-clientChan
}

// Main returns true if this room is a main room.
func (ch *Room) Main() bool {
	return ch.name == main
}

// Name returns chanel's name.
func (ch *Room) Name() string {
	return ch.name
}

// SendToEveryone sends message to everyone from that room.
func (ch *Room) SendToEveryone(msg message.Message) {
	ch.incomingMessages <- msg
}

// AddClient adds client to room.
func (ch *Room) AddClient(client message.Sender) {
	ch.addClientChan <- client
}

// RemoveClient removes client from the room.
func (ch *Room) RemoveClient(clientID string) {
	ch.removeClientChan <- clientID
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
	go room.start()
	return room
}

// NewMainRoom returns new unremovable Room struct with name 'main'.
func NewMainRoom(rooms Rooms) *Room {
	return NewRoom(main, rooms)
}
