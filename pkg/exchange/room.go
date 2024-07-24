package exchange

import (
	"fmt"

	logger "github.com/sirupsen/logrus"
)

const (
	// main is the name of the main room.
	main = "main"
)

// MainRoomName returns name of the main room.
func MainRoomName() string {
	return main
}

// NewRoom functions returns new Room struct.
func NewRoom(name string, rooms *Rooms) *Room {
	return &Room{
		name:             name,
		clients:          map[string]*Client{},
		rooms:            rooms,
		clientExists:     make(chan clientExist, 5),
		removeClientChan: make(chan string, 5),
		addClientChan:    make(chan *Client, 5),
		incomingMessages: make(chan *Message, 50),
		interrupt:        make(chan bool, 5),
	}
}

// NewMainRoom returns new unremovable Room struct with name 'main'.
func NewMainRoom(rooms *Rooms) *Room {
	return NewRoom(main, rooms)
}

// Room represents chat room.
type Room struct {
	name             string
	clients          map[string]*Client
	rooms            *Rooms
	removeClientChan chan string
	addClientChan    chan *Client
	clientExists     chan clientExist
	incomingMessages chan *Message
	interrupt        chan bool
}

// FindClient returns client with given id if it exist in this room.
func (ch *Room) FindClient(clientID string) (*Client, error) {
	clientChan := make(chan *Client, 1)

	ch.clientExists <- clientExist{
		existChan: clientChan,
		clientID:  clientID,
	}

	client := <-clientChan
	if client == nil {
		return nil, fmt.Errorf("client with id %v cannot be found", clientID)
	}

	return client, nil
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
func (ch *Room) SendToEveryone(msg *Message) {
	ch.incomingMessages <- msg
}

// AddClient adds client to this room.
func (ch *Room) AddClient(client *Client) {
	if ch == nil {
		logger.Infof("cannot add client: %v to channel, channel doesn't exist", client.ID())
		return
	}
	ch.addClientChan <- client
}

// RemoveClient removes client from this room.
func (ch *Room) RemoveClient(clientID string) {
	if ch == nil {
		logger.Infof("cannot remove client: %v from channel, channel doesn't exist", clientID)
		return
	}
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
					logger.Infof("Room: '%v' is empty. Should be removed.", ch.Name())
					ch.rooms.RemoveRoom(ch.Name())
					ch.interrupt <- true
				}

			case client := <-ch.addClientChan:
				ch.clients[client.ID()] = client

			case msg := <-ch.incomingMessages:
				logger.Infof("Sending msg to %v room members.", len(ch.clients))
				for _, client := range ch.clients {
					logger.Infof("Sending msg to %v from room '%v'.", client, ch.name)
					client.Send(msg)
				}

			case c := <-ch.clientExists:
				c.existChan <- ch.clients[c.clientID]
			}
		}
	}()
}

type clientExist struct {
	existChan chan *Client
	clientID  string
}
