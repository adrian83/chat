package client

import (
	"github.com/adrian83/chat/chat/logger"

	"fmt"

	"github.com/adrian83/chat/chat/ws"
	"github.com/adrian83/chat/chat/ws/message"
)

// Connection defines interface for connections.
type Connection interface {
	Send(msg interface{}) error
	Receive(msg interface{}) error
	Close() error
}

type User interface {
	Name() string
}

// NewClient returns new Client instance
func NewClient(ID string, user User, rooms *ws.DefaultRooms, conn Connection) *DefaultClient {
	client := DefaultClient{
		user:               user,
		id:                 ID,
		rooms:              rooms,
		messagesChannel:    make(chan message.Message, 50),
		stopSendingChannel: make(chan bool, 5),
		stop:               make(chan bool, 5),
		connnection:        conn,
	}

	return &client
}

// DefaultClient default implementation of the Client interface.
type DefaultClient struct {
	id                 string
	user               User
	rooms              *ws.DefaultRooms
	messagesChannel    chan message.Message
	stopSendingChannel chan bool
	stop               chan bool
	connnection        Connection
}

func (c *DefaultClient) Rooms() {
}

func (c *DefaultClient) Start() {
	c.startSending()
	c.startReceiving()
	<-c.stop
}

// StartSending starts infinite loop which is sending messages.
func (c *DefaultClient) startSending() {
	go func() {
		for {
			select {
			case msg := <-c.messagesChannel:
				logger.Infof("Client", "StartSending", "Client %v will send message %v", c, msg)

				if err := c.connnection.Send(msg); err != nil {
					logger.Warnf("Client", "StartSending", "Client %v cannot send message.Error: %v", c, err)
				}

			case b := <-c.stopSendingChannel:
				logger.Infof("Client", "StartSending", "Client %v, stop %v", c, b)
				c.rooms.RemoveClient(c)
				c.stop <- true
				return
			}
		}
	}()
}

// StartReceiving starts infinite loop which is processing received messages.
func (c *DefaultClient) startReceiving() {

	go func() {
		defer c.closeConnection()

		for {

			msg := new(message.Message)
			if err := c.connnection.Receive(msg); err != nil {
				logger.Infof("Client", "StartReceiving", "Error in Client %v while reading from websocket. Error: %v", c, err)
				c.stopSendingChannel <- true
				return
			}

			msg.SenderName = c.user.Name()
			msg.SenderID = c.id

			logger.Infof("Client", "StartReceiving", "Client %v received a messanges: %v", c, msg)

			msg.DoWith(c, c.rooms)

		}
	}()
}

// ID returns id of the client.
func (c *DefaultClient) ID() string {
	return c.id
}

// String returns string representation of DefaultClient struct.
func (c *DefaultClient) String() string {
	return fmt.Sprintf("Client { name: %v }", c.user.Name())
}

// Send sends message through connection.
func (c *DefaultClient) Send(msg message.Message) {
	c.messagesChannel <- msg
}

func (c *DefaultClient) closeConnection() {
	if err := c.connnection.Close(); err != nil {
		logger.Warnf("Client", "closeConnection", "Error in %v while closing connection. Error: %v", c, err)
	}
}
