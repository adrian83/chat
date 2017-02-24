package ws

import (
	"db"
	"logger"

	"fmt"

	"golang.org/x/net/websocket"
)

// NewClient returns new Client instance
func NewClient(ID string, user db.User, wsc *websocket.Conn, allChannels *Channels) *Client {
	client := Client{
		connection:  NewConnection(wsc),
		interrupt:   make(chan bool, 5),
		user:        user,
		id:          ID,
		allChannels: allChannels,
	}

	return &client
}

type Client struct {
	id          string
	user        db.User
	connection  *WsConnection
	interrupt   chan bool
	allChannels *Channels
}

func (c *Client) String() string {
	return fmt.Sprintf("Client { name: %v }", c.user.Name)
}

func (c *Client) Send(msg Message) error {
	return c.connection.Send(msg)
}

func (c *Client) Start() {

	channelNamesMsg := Message{
		MsgType:    "CHAN_LIST_MSG",
		SenderID:   "system",
		SenderName: "system",
		Channels:   c.allChannels.Names(),
	}

	err := c.connection.Send(channelNamesMsg)
	if err != nil {
		logger.Infof("Client", "Start", "Error while sending message: %v", err)
	}

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

			msgOld := msgToClient.Message
			msg := Message{
				MsgType:    msgOld.MsgType,
				SenderID:   msgOld.SenderID,
				SenderName: c.user.Name,
				Channels:   msgOld.Channels,
				Channel:    msgOld.Channel,
				Content:    msgOld.Content,
			}

			switch msg.MsgType {
			case "LOGOUT_USER":
				logger.Infof("Client", "Start", "%v received incomming message: %v", c, msg.MsgType)

				// remove client from channels
				for _, channel := range c.allChannels.ClientsChannels(c) {
					if errors := channel.RemoveClient(c); len(errors) > 0 {
						for _, sendErr := range errors {
							logger.Warnf("Client", "Start", "Error while sending 'TEXT_MSG' message. Error: %v", sendErr)
						}
					}
				}

				// stop client
				c.interrupt <- true

				// TODO stop connection

				// resend message
				if err = c.connection.Send(msgToClient.Message); err != nil {
					logger.Infof("Client", "Start", "Error while sending message: %v", err)
				}

			case "TEXT_MSG":

				logger.Warnf("Client", "Start", "Received 'TEXT_MSG' message on channel %v.", msg.Channel)
				channel, ok := c.allChannels.ClientsChannels(c)[msg.Channel]
				if ok {
					if errors := channel.SendToEveryone(msg); len(errors) > 0 {
						for _, sendErr := range errors {
							logger.Warnf("Client", "Start", "Error while sending 'TEXT_MSG' message. Error: %v", sendErr)
						}
					}

				}
			case "ADD_CH":

				_, ok := c.allChannels.ClientsChannels(c)[msg.Channel]
				if ok {
					logger.Infof("Client", "Start", "Channel: %v already exists", msg.Channel)
				} else {
					logger.Infof("Client", "Start", "Channel: %v does not exist", msg.Channel)
					channel := NewChannel(msg.Channel, c, c.allChannels)
					c.allChannels.AddChannel(channel)
				}

			case "USER_JOINED_CH":

				if _, ok := c.allChannels.ClientsChannels(c)[msg.Channel]; !ok {
					c.allChannels.AddClientToChannel(msg.Channel, c)
					if err = c.connection.Send(msgToClient.Message); err != nil {
						logger.Infof("Client", "Start", "Error while sending message: %v", err)
					}
				}

			case "USER_LEFT_CH":

				if errors := c.allChannels.RemoveClientFromChannel(msg.Channel, c); len(errors) > 0 {
					for _, sendErr := range errors {
						logger.Warnf("Client", "Start", "Error while sending 'USER_LEFT_CH' message. Error: %v", sendErr)
					}
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

// Stop stops client goroutine
func (c *Client) Stop() {
	c.interrupt <- true
	c.connection.Close()
}
