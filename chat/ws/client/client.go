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
		user:               user,
		id:                 ID,
		rooms:              rooms,
		messagesChannel:    make(chan message.Message, 50),
		stopSendingChannel: make(chan bool, 5),
		stop:               make(chan bool, 5),
		connnection:        conn,
	}
}

// Client represents user of this application.
type Client struct {
	id                 string
	user               User
	rooms              *ws.DefaultRooms
	messagesChannel    chan message.Message
	stopSendingChannel chan bool
	stop               chan bool
	connnection        Connection
}

// Start starts two goroutines: one for sending and one for receiving messages.
func (c *Client) Start() {
	c.startSending()
	c.startReceiving()
	<-c.stop
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
	c.messagesChannel <- msg
}

func (c *Client) closeConnection() {
	if err := c.connnection.Close(); err != nil {
		logger.Warnf("Error in %v while closing connection. Error: %v", c, err)
	}
}

// StartSending starts infinite loop which is sending messages.
func (c *Client) startSending() {
	go func() {
		for {
			select {
			case msg := <-c.messagesChannel:
				logger.Infof("Client %v will send message %v", c, msg)

				if err := c.connnection.Send(msg); err != nil {
					logger.Warnf("Client %v cannot send message.Error: %v", c, err)
				}

			case b := <-c.stopSendingChannel:
				logger.Infof("Client %v, stop %v", c, b)
				c.rooms.RemoveClient(c)
				c.stop <- true
				return
			}
		}
	}()
}

// startReceiving starts infinite loop which is processing received messages.
func (c *Client) startReceiving() {

	go func() {
		defer c.closeConnection()

		for {

			msg := new(message.Message)
			if err := c.connnection.Receive(msg); err != nil {
				logger.Infof("Error in Client %v while reading from websocket. Error: %v", c, err)
				c.stopSendingChannel <- true
				return
			}

			msg.SenderName = c.user.Name()
			msg.SenderID = c.id

			logger.Infof("Client %v received a messanges: %v", c, msg)

			msg.DoWith(c, c.rooms)

		}
	}()
}
