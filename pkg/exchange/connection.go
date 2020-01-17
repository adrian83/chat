package exchange

import (
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

// NewWebSocketConn returns new instance of wsConnection,
func NewWebSocketConn(webSocketConn *websocket.Conn) *WsConnection {
	return &WsConnection{
		webSocketConn: webSocketConn,
	}
}

type WsConnection struct {
	webSocketConn *websocket.Conn
}

func (c *WsConnection) Send(msg interface{}) error {
	err := websocket.JSON.Send(c.webSocketConn, msg)
	return errors.Wrapf(err, "error while sending message through websocket")
}

func (c *WsConnection) Receive(msg interface{}) error {
	err := websocket.JSON.Receive(c.webSocketConn, &msg)
	return errors.Wrapf(err, "error while receiving message from websocket")
}

func (c *WsConnection) Close() error {
	err := c.webSocketConn.Close()
	return errors.Wrapf(err, "error while closing websocket connection")
}
