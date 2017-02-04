package ws

import (
	"db"
	"logger"

	"fmt"

	"golang.org/x/net/websocket"
)

func NewClient(ID string, user db.User, wsc *websocket.Conn) *Client {
	client := Client{
		connection: NewConnection(wsc),
		interrupt:  make(chan bool),
		user:       user,
		id:         ID,
	}

	return &client
}

type Client struct {
	id         string
	user       db.User
	connection *WsConnection
	interrupt  chan bool
}

func (c *Client) String() string {
	return fmt.Sprintf("Client { name: %v }", c.user.Name)
}

func (c *Client) Start() {

	msg := map[string]interface{}{
		"msgType":    "CHAN_LIST_MSG",
		"senderId":   "0",
		"senderName": "system",
		"receiver":   c.id,
		"channels":   []string{"abc", "main", "dfg"},
	}

	err := c.connection.Send(msg)
	if err != nil {
		logger.Infof("Client", "Start", "Error while sending message: %v", err)
	}

outer:
	for {
		logger.Infof("Client", "Start", "%v waithing for a msg", c)
		select {
		case msgToClient := <-c.connection.ToClient:
			logger.Infof("Client", "Start", "%v received incomming message: %v", c, msgToClient)
			//c.server.Receive(incommingMsg)

		case exit := <-c.interrupt:
			if exit {
				logger.Infof("Client", "Start", "%v received interupt msg", c)
				close(c.interrupt)
				break outer
			}
		}
	}
	logger.Infof("Client", "Start", "%v end", c)
}

// Stop stops client goroutine
func (c *Client) Stop() {
	c.interrupt <- true
	c.connection.Close()
}
