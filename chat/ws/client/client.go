package client

import (
	"fmt"

	"github.com/adrian83/chat/chat/ws"
	"github.com/adrian83/chat/chat/ws/message"

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
func NewClient(ID string, user User, rooms *ws.DefaultRooms, conn Connection) *Client {
	return &Client{
		user:        user,
		id:          ID,
		rooms:       rooms,
		messages:    make(chan message.Message, 50),
		stopSending: make(chan interface{}, 1),
		stop:        make(chan interface{}, 1),
		connnection: conn,
	}
}

// Client represents user of this application.
type Client struct {
	id          string
	user        User
	rooms       *ws.DefaultRooms
	messages    chan message.Message
	stopSending chan interface{}
	stop        chan interface{}
	connnection Connection
}

// Start starts two goroutines: one for sending and one for receiving messages.
func (c *Client) Start() {
	logger.Infof("Client: %v. Starting", c.user.Name())

	defer c.closeConnection()

	c.startSending()
	c.startReceiving()

	<-c.stop

	logger.Infof("Client: %v. Stoping", c.user.Name())
}

// ID returns id of the client.
func (c *Client) ID() string {
	return c.id
}

// String is a string representation of Client struct.
func (c *Client) String() string {
	return fmt.Sprintf("{\"name\":\"%v\"}", c.user.Name())
}

// Send sends message through connection.
func (c *Client) Send(msg message.Message) {
	logger.Infof("Client: %v. Adding message to send channel. Message: %v", c.user.Name(), msg.MsgType)
	c.messages <- msg
}

func (c *Client) closeConnection() {
	logger.Infof("Client: %v. Closing connection", c.user.Name())
	if err := c.connnection.Close(); err != nil {
		logger.Warnf("Client: %v. Error while closing connection. Error: %v", c.user.Name(), err)
	}
}

// StartSending starts infinite loop which is sending messages.
func (c *Client) startSending() {
	logger.Infof("Client: %v. Starting sending messages", c.user.Name())
	go func() {
		for {
			select {
			case msg := <-c.messages:
				logger.Infof("Client: %v. Sending message. Message: %v", c.user.Name(), msg.MsgType)

				if err := c.connnection.Send(msg); err != nil {
					logger.Warnf("Client: %v. Error while sending message.Error: %v", c.user.Name(), err)
				}

			case <-c.stopSending:
				logger.Infof("Client: %v. Stopping sending messages", c.user.Name())
				c.rooms.RemoveClient(c)
				c.stop <- true
				return
			}
		}
	}()
	logger.Infof("Client: %v. Stopping sending messages", c.user.Name())
}

// startReceiving starts infinite loop which is processing received messages.
func (c *Client) startReceiving() {
	logger.Infof("Client: %v. Starting receiving messages", c.user.Name())
	go func() {

		for {

			msg := new(message.Message)
			if err := c.connnection.Receive(msg); err != nil {
				logger.Warnf("Client: %v. Error while receiving message. Error: %v", c.user.Name(), err)
				c.stopSending <- true
				return
			}

			msg.SenderName = c.user.Name()
			msg.SenderID = c.id

			logger.Infof("Client: %v. Received message. Message: %v", c.user.Name(), msg.MsgType)

			msg.DoWith(c, c.rooms)
		}
	}()
	logger.Infof("Client: %v. Stopping receiving messages", c.user.Name())
}
