package ws

import (
	"regexp"

	"chat/db"
	"chat/logger"

	"fmt"

	"golang.org/x/net/websocket"
)

var (
	channelNameRegexp = `^[a-zA-Z0-9_.-]*$`
	validChannelName  = regexp.MustCompile(channelNameRegexp)
)

// NewClient returns new Client instance
func NewClient(ID string, user db.User, channels Channels, wsConnnection *websocket.Conn) Client {
	client := DefaultClient{
		user:     user,
		id:       ID,
		channels: channels,

		messagesToSend: make(chan Message, 50),

		stopSending: make(chan bool, 5),

		wsConnnection: wsConnnection,
	}

	return &client
}

// Client interface definig client of the app.
type Client interface {
	Send(msg Message)
	ID() string

	StartSending()
	StartReceiving()
}

// DefaultClient default implementation of the Client interface.
type DefaultClient struct {
	id   string
	user db.User

	channels Channels

	messagesToSend chan Message
	wsConnnection  *websocket.Conn
	stopSending    chan bool
}

// StartSending starts infinite loop which is sending messages.
func (c *DefaultClient) StartSending() {
	logger.Infof("Client", "StartSending", "Client %v is starting sending messanges", c)
mainLoop:
	for {
		select {
		case s := <-c.messagesToSend:
			logger.Infof("Client", "StartSending", "Client %v, message %v", c, s)
		case b := <-c.stopSending:
			logger.Infof("Client", "StartSending", "Client %v, stop %v", c, b)
			c.channels.RemoveClient(c)
			break mainLoop
		}
	}
	logger.Infof("Client", "StartSending", "Client %v is stopping sending messanges", c)
}

// StartReceiving starts infinite loop which is processing received messages.
func (c *DefaultClient) StartReceiving() {
	logger.Infof("Client", "StartReceiving", "Client %v is starting receiving messanges", c)

	defer func() {
		if err := c.wsConnnection.Close(); err != nil {
			logger.Warnf("Client", "StartReceiving", "Error in %v while closing connection. Error: %v", c, err)
		}
	}()

mainLoop:
	for {

		msg := Message{}
		if err := websocket.JSON.Receive(c.wsConnnection, &msg); err != nil {
			logger.Infof("Client", "StartReceiving", "Error in Client %v while reading from websocket. Error: %v", c, err)
			c.stopSending <- true
			break mainLoop
		}

		msg.SenderName = c.user.Name
		msg.SenderID = c.id

		c.handleMessage(msg)

	}
	logger.Infof("Client", "StartReceiving", "Client %v is stopping receiving messanges", c)
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
func (c *DefaultClient) Send(msg Message) {
	c.messagesToSend <- msg
}

func (c *DefaultClient) handleMessage(msg Message) {
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
}

func (c *DefaultClient) logSendErrors(msg Message, errors []SendError) {
	for _, sendErr := range errors {
		logger.Warnf("Client", "logSendErrors", "Error while sending message %v to %v. Error: %v", msg, sendErr.Client, sendErr.Err)
	}
}

func (c *DefaultClient) handleLogoutMessage(msg Message) {
	// remove client from channels
	logger.Infof("Client", "Start", "Logging out client: %v", c)

	// resend message
	c.Send(msg)

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

	if !validChannelName.MatchString(msg.Channel) {
		errMsg := ErrorMessage(fmt.Sprintf("Invalid channel name. Name must match %v", channelNameRegexp))
		c.Send(errMsg)
		return
	}

	channel := NewChannel(msg.Channel, c, c.channels)
	if errors := c.channels.AddChannel(channel); len(errors) > 0 {
		c.logSendErrors(msg, errors)
	}
}

func (c *DefaultClient) handleJoinChannelMessage(msg Message) {
	if _, ok := c.channels.ClientsChannels(c)[msg.Channel]; !ok {
		c.channels.AddClientToChannel(msg.Channel, c)
		c.Send(msg)
	}
}

func (c *DefaultClient) handleLeaveChannelMessage(msg Message) {
	if errors := c.channels.RemoveClientFromChannel(msg.Channel, c); len(errors) > 0 {
		c.logSendErrors(msg, errors)
	}
}
