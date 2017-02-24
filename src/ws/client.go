package ws

import (
	"db"
	"logger"

	"fmt"

	"golang.org/x/net/websocket"
)

// NewClient returns new Client instance
func NewClient(ID string, user db.User, wsc *websocket.Conn, channels *Channels) *Client {
	client := Client{
		connection: NewConnection(wsc),
		interrupt:  make(chan bool, 5),
		user:       user,
		id:         ID,
		channels:   channels,
	}

	return &client
}

type Client struct {
	id         string
	user       db.User
	connection *WsConnection
	interrupt  chan bool
	channels   *Channels
}

func (c *Client) String() string {
	return fmt.Sprintf("Client { name: %v }", c.user.Name)
}

func (c *Client) Send(msg Message) error {
	return c.connection.Send(msg)
}

func (c *Client) Start() {

outer:
	for {
		logger.Infof("Client", "Start", "%v waithing for a msg", c)
		select {
		case msgToClient := <-c.connection.ToClient:
			logger.Infof("Client", "Start", "%v received incomming message: %v", c, msgToClient)

			if msgToClient.Error != nil {
				logger.Warnf("Client", "Start", "Error while getting message for client %v. Error: %v", c, msgToClient.Error)
				continue
			}

			msg := msgToClient.Message
			msg.SenderName = c.user.Name
			msg.SenderID = c.id

			switch msg.MsgType {
			case "LOGOUT_USER":

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
				if err := c.Send(msgToClient.Message); err != nil {
					logger.Infof("Client", "Start", "Error while sending message: %v", err)
				}

			case "TEXT_MSG":

				channel, ok := c.channels.ClientsChannels(c)[msg.Channel]
				if ok {
					if errors := channel.SendToEveryone(msg); len(errors) > 0 {
						c.logSendErrors(msg, errors)
					}

				}
			case "ADD_CH":

				channel := NewChannel(msg.Channel, c, c.channels)
				if errors := c.channels.AddChannel(channel); len(errors) > 0 {
					c.logSendErrors(msg, errors)
				}

			case "USER_JOINED_CH":

				if _, ok := c.channels.ClientsChannels(c)[msg.Channel]; !ok {
					c.channels.AddClientToChannel(msg.Channel, c)
					if err := c.connection.Send(msgToClient.Message); err != nil {
						logger.Warnf("Client", "Start", "Error while sending message %v to %v. Error: %v", msg, c, err)
					}
				}

			case "USER_LEFT_CH":
				if errors := c.channels.RemoveClientFromChannel(msg.Channel, c); len(errors) > 0 {
					c.logSendErrors(msg, errors)
				}

			default:
				logger.Infof("Client", "Start", "Unknown message: %v", msg.MsgType)
			}

			//c.server.Receive(incommingMsg)

		case exit := <-c.interrupt:
			logger.Info("Client", "Start", "Interrupted")
			if exit {
				logger.Infof("Client", "Start", "%v received interupt msg", c)
				close(c.interrupt)
				break outer
			}
		}
	}
	logger.Infof("Client", "Start", "%v end", c)
}

func (c *Client) logSendErrors(msg Message, errors []SendError) {
	for _, sendErr := range errors {
		logger.Warnf("Client", "logSendErrors", "Error while sending message %v to %v. Error: %v", msg, sendErr.Client, sendErr.Err)
	}
}

// Stop stops client goroutine
func (c *Client) Stop() {
	c.interrupt <- true
	c.connection.Close()
}
