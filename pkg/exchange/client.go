package exchange

import (
	"fmt"

	logger "github.com/sirupsen/logrus"
)

// User is an interface whichh defines persisten data about application user.
type user interface {
	Name() string
}

// NewClient returns new Client instance
func NewClient(id string, user user, rooms *Rooms, conn *WsConnection, router *Router) *Client {
	return &Client{
		user:        user,
		id:          id,
		rooms:       rooms,
		connnection: conn,
		router:      router,
		messages:    make(chan *Message, 50),
		stopSending: make(chan interface{}, 1),
		stopWaiting: make(chan interface{}, 1),
	}
}

// Client represents user of this application.
type Client struct {
	id          string
	user        user
	rooms       *Rooms
	router      *Router
	connnection *WsConnection
	messages    chan *Message
	stopSending chan interface{}
	stopWaiting chan interface{}
}

// Start starts two goroutines: one for sending and one for receiving messages.
func (c *Client) Start() {
	logger.Infof("Client: %v. Starting", c.user.Name())

	defer c.closeConnection()

	c.startSending()
	c.startReceiving()

	<-c.stopWaiting

	logger.Infof("Client: %v. Stoping", c.user.Name())
}

// ID returns id of the client.
func (c *Client) ID() string {
	return c.id
}

// String is a string representation of Client struct.
func (c *Client) String() string {
	return fmt.Sprintf(`{"name":"%v"}`, c.user.Name())
}

// Send sends message through connection.
func (c *Client) Send(msg *Message) {
	logger.Infof("Client: %v. Adding message to send channel. Message: %v", c.user.Name(), msg.MsgType)
	c.messages <- msg
}

func (c *Client) closeConnection() {
	logger.Infof("Client: %v. Closing connection", c.user.Name())
	if err := c.connnection.Close(); err != nil {
		logger.Warnf("Client: %v. Error while closing connection. Error: %v", c.user.Name(), err)
	}
}

func (c *Client) stop() {
	c.stopSending <- true
	c.stopWaiting <- true
}

// StartSending starts infinite loop which is sending messages.
func (c *Client) startSending() {
	logger.Infof("Client: %v. Starting sending messages", c.user.Name())
	go func() {
	mainLoop:
		for {
			select {
			case msg := <-c.messages:
				logger.Infof("Client: %v. Sending message. Message: %v", c.user.Name(), msg.MsgType)

				if err := c.connnection.Send(msg); err != nil {
					logger.Warnf("Client: %v. Error while sending message.Error: %v", c.user.Name(), err)
					c.stop()
				}

			case <-c.stopSending:
				logger.Infof("Client: %v. Stopping sending messages", c.user.Name())
				c.rooms.RemoveClient(c)
				break mainLoop
			}
		}
		logger.Infof("Client: %v. Stopping sending messages", c.user.Name())
	}()

}

// startReceiving starts infinite loop which is processing received messages.
func (c *Client) startReceiving() {
	logger.Infof("Client: %v. Starting receiving messages", c.user.Name())
	go func() {

		for {
			var msg Message
			if err := c.connnection.Receive(&msg); err != nil {
				logger.Warnf("Client: %v. Error while receiving message. Error: %v", c.user.Name(), err)
				c.stop()
				break
			}

			msg.SenderName = c.user.Name()
			msg.SenderID = c.id

			logger.Infof("Client: %v. Received message. Message: %v", c.user.Name(), msg.MsgType)

			if err := c.router.FindRoute(msg.MsgType).Handle(&msg); err != nil {
				logger.Warnf("Client: %v. Error while handling message: %v. Error: %v", c.user.Name(), msg, err)
			}
		}

		logger.Infof("Client: %v. Stopping receiving messages", c.user.Name())
	}()

}
