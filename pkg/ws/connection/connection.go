package connection

import (
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

// NewWebSocketConn returns new instance of wsConnection,
func NewWebSocketConn(webSocketConn *websocket.Conn) *wsConnection {
	return &wsConnection{
		webSocketConn: webSocketConn,
	}
}

type wsConnection struct {
	webSocketConn *websocket.Conn
}

func (c *wsConnection) Send(msg interface{}) error {
	err := websocket.JSON.Send(c.webSocketConn, msg)
	return errors.Wrapf(err, "error while sending message through websocket")
}

func (c *wsConnection) Receive(msg interface{}) error {
	err := websocket.JSON.Receive(c.webSocketConn, &msg)
	return errors.Wrapf(err, "error while receiving message from websocket")
}

func (c *wsConnection) Close() error {
	err := c.webSocketConn.Close()
	return errors.Wrapf(err, "error while closing websocket connection")
}
