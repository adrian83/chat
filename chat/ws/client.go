package ws

import (
	"github.com/adrian83/go-chat/chat/db"
	"github.com/adrian83/go-chat/chat/logger"

	"fmt"

	"golang.org/x/net/websocket"
)

// NewClient returns new Client instance
func NewClient(ID string, user db.User, rooms Rooms, wsConnnection *websocket.Conn) Client {
	client := DefaultClient{
		user:               user,
		id:                 ID,
		rooms:              rooms,
		messagesChannel:    make(chan Message, 50),
		stopSendingChannel: make(chan bool, 5),
		wsConnnection:      wsConnnection,
	}

	return &client
}

// Client interface definig client of the app.
type Client interface {
	ID() string
	Send(msg Message)
	StartSending()
	StartReceiving()
}

// DefaultClient default implementation of the Client interface.
type DefaultClient struct {
	id                 string
	user               db.User
	rooms              Rooms
	messagesChannel    chan Message
	stopSendingChannel chan bool
	wsConnnection      *websocket.Conn
}

// StartSending starts infinite loop which is sending messages.
func (c *DefaultClient) StartSending() {
	logger.Infof("Client", "StartSending", "Client %v is starting sending messanges", c)
mainLoop:
	for {
		select {
		case msg := <-c.messagesChannel:
			logger.Infof("Client", "StartSending", "Client %v will send message %v", c, msg)

			if err := websocket.JSON.Send(c.wsConnnection, msg); err != nil {
				logger.Warnf("Client", "StartSending", "Client %v cannot send message.Error: %v", c, err)
			}

		case b := <-c.stopSendingChannel:
			logger.Infof("Client", "StartSending", "Client %v, stop %v", c, b)
			c.rooms.RemoveClient(c)
			break mainLoop
		}
	}
	logger.Infof("Client", "StartSending", "Client %v is stopping sending messanges", c)
}

// StartReceiving starts infinite loop which is processing received messages.
func (c *DefaultClient) StartReceiving() {
	logger.Infof("Client", "StartReceiving", "Client %v is starting receiving messanges", c)

	defer c.closeConnection()

mainLoop:
	for {

		msg := Message{}
		if err := websocket.JSON.Receive(c.wsConnnection, &msg); err != nil {
			logger.Infof("Client", "StartReceiving", "Error in Client %v while reading from websocket. Error: %v", c, err)
			c.stopSendingChannel <- true
			break mainLoop
		}

		msg.SenderName = c.user.Name
		msg.SenderID = c.id

		logger.Infof("Client", "StartReceiving", "Client %v received a messanges: %v", c, msg)

		switch msg.MsgType {
		case LogoutMT:
			c.Send(msg)
		case TextMsgMT:
			c.rooms.SendMessageOnRoom(msg)
		case CreateRoomMT:
			c.rooms.CreateRoom(msg.Room, c)
		case UserJoinedRoomMT:
			c.rooms.AddClientToRoom(msg.Room, c)
		case UserLeftRoomMT:
			c.rooms.RemoveClientFromRoom(msg.Room, c)
		default:
			logger.Infof("Client", "Start", "Unknown message: %v", msg.MsgType)
		}

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
	c.messagesChannel <- msg
}

func (c *DefaultClient) closeConnection() {
	if err := c.wsConnnection.Close(); err != nil {
		logger.Warnf("Client", "closeConnection", "Error in %v while closing connection. Error: %v", c, err)
	}
}
