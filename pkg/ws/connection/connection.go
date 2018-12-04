package connection

import "golang.org/x/net/websocket"

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
	return websocket.JSON.Send(c.webSocketConn, msg)
}

func (c *wsConnection) Receive(msg interface{}) error {
	return websocket.JSON.Receive(c.webSocketConn, &msg)
}

func (c *wsConnection) Close() error {
	return c.webSocketConn.Close()
}
