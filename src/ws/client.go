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

func (c *Client) Send(msg map[string]interface{}) error {
	return c.connection.Send(msg)
}

func (c *Client) Start() {

	msg := map[string]interface{}{
		"msgType":    "CHAN_LIST_MSG",
		"senderId":   "0",
		"senderName": "system",
		"receiver":   c.id,
		"channels":   c.allChannels.Names(),
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

			if msgToClient.Error != nil {
				logger.Warnf("Client", "Start", "Error while getting message for client %v. Error: %v", c, msgToClient.Error)
				continue
			}

			msg := Message(msgToClient.Message)
			msgType, err := msg.getType()
			if err != nil {
				logger.Warnf("Client", "Start", "Error while getting type of message. Error: %v", err)
				continue
			}

			switch msgType {
			case "LOGOUT_USER":
				logger.Infof("Client", "Start", "%v received incomming message: %v", c, msgType)

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
				msgToClient.Message["senderName"] = c.user.Name
				channelName := msgToClient.Message["channel"]
				logger.Warnf("Client", "Start", "Received 'TEXT_MSG' message on channel %v.", channelName.(string))
				channel, ok := c.allChannels.ClientsChannels(c)[channelName.(string)]
				if ok {
					if errors := channel.SendToEveryone(msgToClient.Message); len(errors) > 0 {
						for _, sendErr := range errors {
							logger.Warnf("Client", "Start", "Error while sending 'TEXT_MSG' message. Error: %v", sendErr)
						}
					}

				}
			case "ADD_CH":
				channelName := msgToClient.Message["channel"]
				_, ok := c.allChannels.ClientsChannels(c)[channelName.(string)]
				if ok {
					logger.Infof("Client", "Start", "Channel: %v already exists", channelName)
				} else {
					logger.Infof("Client", "Start", "Channel: %v does not exist", channelName)
					channel := NewChannel(channelName.(string), c, c.allChannels)
					c.allChannels.AddChannel(channel)
				}

			case "USER_JOINED_CH":
				channelName := msgToClient.Message["channel"]
				if _, ok := c.allChannels.ClientsChannels(c)[channelName.(string)]; !ok {
					c.allChannels.AddClientToChannel(channelName.(string), c)
					if err = c.connection.Send(msgToClient.Message); err != nil {
						logger.Infof("Client", "Start", "Error while sending message: %v", err)
					}
				}

			default:
				logger.Infof("Client", "Start", "Unknown message: %v", msgType)
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
