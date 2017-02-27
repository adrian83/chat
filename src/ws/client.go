package ws

import (
	"db"
	"logger"

	"fmt"
	"io"
)

// NewClient returns new Client instance
func NewClient(ID string, user db.User, conn Connection, channels *Channels) Client {
	client := DefaultClient{
		connection: conn,
		interrupt:  make(chan bool, 5),
		user:       user,
		id:         ID,
		channels:   channels,
	}

	return &client
}

// Client interface definig client of the app.
type Client interface {
	Send(msg Message) error
	Start()
	Stop()
	ID() string
}

// DefaultClient default implementation of the Client interface.
type DefaultClient struct {
	id         string
	user       db.User
	connection Connection
	interrupt  chan bool
	channels   *Channels
}

// ID returns id of the client.
func (c *DefaultClient) ID() string {
	return c.id
}

// String returns string representation of DefaultClient struct.
func (c *DefaultClient) String() string {
	return fmt.Sprintf("Client { name: %v }", c.user.Name)
}

// Send sends message through connection.
func (c *DefaultClient) Send(msg Message) error {
	return c.connection.Send(msg)
}

// Start starts client main loop.
func (c *DefaultClient) Start() {

outer:
	for {
		logger.Infof("Client", "Start", "%v waithing for a msg", c)
		select {
		case msgToClient := <-c.connection.Incomming():
			logger.Infof("Client", "Start", "%v received incomming message: %v", c, msgToClient)

			msg := msgToClient.Message
			msg.SenderName = c.user.Name
			msg.SenderID = c.id

			if msgToClient.Error != nil {
				if msgToClient.Error == io.EOF {
					c.handleLogoutMessage(msg)
					break outer
				}

				logger.Warnf("Client", "Start", "Error while getting message for client %v. Error: %v", c, msgToClient.Error)
				continue
			}

			switch msg.MsgType {
			case "LOGOUT_USER":
				c.handleLogoutMessage(msg)
			case "TEXT_MSG":
				c.handleTextMessage(msg)
			case "ADD_CH":
				c.handleCreateChannelMessage(msg)
			case "USER_JOINED_CH":
				c.handleJoinChannelMessage(msg)
			case "USER_LEFT_CH":
				c.handleLeaveChannelMessage(msg)
			default:
				logger.Infof("Client", "Start", "Unknown message: %v", msg.MsgType)
			}

		case <-c.interrupt:
			logger.Infof("Client", "Start", "%v received interupt msg", c)
			close(c.interrupt)
			break outer

		}
	}
	logger.Infof("Client", "Start", "%v end", c)
}

func (c *DefaultClient) logSendErrors(msg Message, errors []SendError) {
	for _, sendErr := range errors {
		logger.Warnf("Client", "logSendErrors", "Error while sending message %v to %v. Error: %v", msg, sendErr.Client, sendErr.Err)
	}
}

// Stop stops client.
func (c *DefaultClient) Stop() {
	c.interrupt <- true
	c.connection.Close()
}

func (c *DefaultClient) handleLogoutMessage(msg Message) {
	// remove client from channels
	for _, channel := range c.channels.ClientsChannels(c) {
		if errors := channel.RemoveClient(c); len(errors) > 0 {
			c.logSendErrors(msg, errors)
		}
	}

	// stop client
	c.interrupt <- true

	// TODO stop connection

	// resend message
	if err := c.Send(msg); err != nil {
		logger.Warnf("Client", "Start", "Error while sending message: %v", err)
	}
}

func (c *DefaultClient) handleTextMessage(msg Message) {
	channel, ok := c.channels.ClientsChannels(c)[msg.Channel]
	if ok {
		if errors := channel.SendToEveryone(msg); len(errors) > 0 {
			c.logSendErrors(msg, errors)
		}
	}
}

func (c *DefaultClient) handleCreateChannelMessage(msg Message) {
	channel := NewChannel(msg.Channel, c, c.channels)
	if errors := c.channels.AddChannel(channel); len(errors) > 0 {
		c.logSendErrors(msg, errors)
	}
}

func (c *DefaultClient) handleJoinChannelMessage(msg Message) {
	if _, ok := c.channels.ClientsChannels(c)[msg.Channel]; !ok {
		c.channels.AddClientToChannel(msg.Channel, c)
		if err := c.connection.Send(msg); err != nil {
			logger.Warnf("Client", "Start", "Error while sending message %v to %v. Error: %v", msg, c, err)
		}
	}
}

func (c *DefaultClient) handleLeaveChannelMessage(msg Message) {
	if errors := c.channels.RemoveClientFromChannel(msg.Channel, c); len(errors) > 0 {
		c.logSendErrors(msg, errors)
	}
}
