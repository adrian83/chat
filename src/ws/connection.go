package ws

import (
	"logger"

	"io"

	"golang.org/x/net/websocket"
)

// Connection interface definig connection which can be used to comunicate between clients.
type Connection interface {
	Send(msg Message) error
	Close() error
	Start()
	Incomming() chan CommunicationResult
}

// CommunicationResult contains result of the  communication (msg or error)
type CommunicationResult struct {
	Message Message
	Error   error
}

// wsConnection represents web socket connection, implements Connection interface.
type wsConnection struct {
	connnection *websocket.Conn
	incomming   chan CommunicationResult
}

// NewConnection returns implementation of Connection interface - in this case wsConnection.
func NewConnection(connnection *websocket.Conn) Connection {
	logger.Info("WsConnection", "NewConnection", "New connection created")
	conn := wsConnection{
		connnection: connnection,
		incomming:   make(chan CommunicationResult, 100),
	}

	go conn.Start()

	return &conn
}

// Start starts the connection.
func (c wsConnection) Start() {
	for {
		logger.Info("wsConnection", "Start", "Waiting for message")

		msg := Message{}
		err := websocket.JSON.Receive(c.connnection, &msg)

		c.incomming <- CommunicationResult{
			Message: msg,
			Error:   err,
		}

		if err == io.EOF {
			logger.Info("wsConnection", "Start", "Received EOF. Exiting")
			break
		}
	}
}

// Incomming returns the channel which can be used to get incomming messages.
func (c wsConnection) Incomming() chan CommunicationResult {
	return c.incomming
}

// Send sends message through the connection.
func (c wsConnection) Send(msg Message) error {
	return websocket.JSON.Send(c.connnection, msg)
}

// Close closes the connection
func (c wsConnection) Close() error {
	return c.connnection.Close()
}
