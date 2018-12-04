package client

import (
	"fmt"

	"github.com/adrian83/chat/pkg/ws/message"
	"github.com/adrian83/chat/pkg/ws/rooms"

	logger "github.com/sirupsen/logrus"
)

// Connection defines interface for connections.
type Connection interface {
	Send(msg interface{}) error
	Receive(msg interface{}) error
	Close() error
}

// User is an interface whichh defines persisten data about application user.
type User interface {
	Name() string
}

// NewClient returns new Client instance
func NewClient(ID string, user User, rooms *rooms.DefaultRooms, conn Connection) *Client {
	return &Client{
		user:        user,
		id:          ID,
		name:        user.Name(),
		rooms:       rooms,
		messages:    make(chan message.Message, 50),
		stopSending: make(chan interface{}, 1),
		stopWaiting: make(chan interface{}, 1),
		connnection: conn,
	}
}

// Client represents user of this application.
type Client struct {
	id          string
	name        string
	user        User
	rooms       *rooms.DefaultRooms
	messages    chan message.Message
	connnection Connection
	stopSending chan interface{}
	stopWaiting chan interface{}
}

// Start starts two goroutines: one for sending and one for receiving messages.
func (c *Client) Start() {
	logger.Infof("Client: %v. Starting", c.name)

	defer c.closeConnection()

	c.startSending()
	c.startReceiving()

	<-c.stopWaiting

	logger.Infof("Client: %v. Stoping", c.name)
}

// ID returns id of the client.
func (c *Client) ID() string {
	return c.id
}

// String is a string representation of Client struct.
func (c *Client) String() string {
	return fmt.Sprintf("{\"name\":\"%v\"}", c.name)
}

// Send sends message through connection.
func (c *Client) Send(msg message.Message) {
	logger.Infof("Client: %v. Adding message to send channel. Message: %v", c.name, msg.MsgType)
	c.messages <- msg
}

func (c *Client) closeConnection() {
	logger.Infof("Client: %v. Closing connection", c.name)
	if err := c.connnection.Close(); err != nil {
		logger.Warnf("Client: %v. Error while closing connection. Error: %v", c.name, err)
	}
}

func (c *Client) stop() {
	c.stopSending <- true
	c.stopWaiting <- true
}

// StartSending starts infinite loop which is sending messages.
func (c *Client) startSending() {
	logger.Infof("Client: %v. Starting sending messages", c.name)
	go func() {
	mainLoop:
		for {
			select {
			case msg := <-c.messages:
				logger.Infof("Client: %v. Sending message. Message: %v", c.name, msg.MsgType)

				if err := c.connnection.Send(msg); err != nil {
					logger.Warnf("Client: %v. Error while sending message.Error: %v", c.name, err)
					c.stop()
				}

			case <-c.stopSending:
				logger.Infof("Client: %v. Stopping sending messages", c.name)
				c.rooms.RemoveClient(c)
				break mainLoop
			}
		}
		logger.Infof("Client: %v. Stopping sending messages", c.name)
	}()

}

// startReceiving starts infinite loop which is processing received messages.
func (c *Client) startReceiving() {
	logger.Infof("Client: %v. Starting receiving messages", c.name)
	go func() {

		for {
			msg := new(message.Message)
			if err := c.connnection.Receive(msg); err != nil {
				logger.Warnf("Client: %v. Error while receiving message. Error: %v", c.name, err)
				c.stop()
				break
			}

			msg.SenderName = c.name
			msg.SenderID = c.id

			logger.Infof("Client: %v. Received message. Message: %v", c.name, msg.MsgType)

			msg.DoWith(c, c.rooms)
		}
		logger.Infof("Client: %v. Stopping receiving messages", c.name)
	}()

}
