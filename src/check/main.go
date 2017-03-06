package main

import (
	"db"
	"logger"
	"time"
	"ws"

	"fmt"
	"math/rand"
	"strconv"
)

func main() {
	fmt.Println("test")

	channels := ws.NewChannels()

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second * 2)
		timeout <- true
	}()

	clients := 30

	for i := 0; i < clients; i++ {

		go func(i int, channels *ws.Channels) {

			strID := strconv.Itoa(i)

			user := db.User{
				ID:       "ID-" + strID,
				Name:     "Name-" + strID,
				Password: "Password-" + strID,
			}

			connection := NewMockConnection()
			go connection.Start()

			client := ws.NewClient(strID, user, connection, channels)
			msg2 := ws.Message{
				MsgType:    "TEXT_MSG",
				SenderID:   "x",
				SenderName: "y",
				//Channels   []string `json:"channels"`
				Channel: "main",
				Content: fmt.Sprintf("MSG { clientID: %v, msgID: %v }", strID, -1),
			}
			connection.Receive(msg2)

			channels.RegisterClient(client)
			time.Sleep(time.Second * 1)
			go client.Start()
			time.Sleep(time.Second * 1)

			newChan := ws.NewChannel("channel-"+strID, client, channels)
			channels.AddChannel(newChan)
			time.Sleep(time.Second * 1)

			time.Sleep(time.Second * 1)
			go func(connection *MockConnection, clientId string, channels *ws.Channels) {
				for ii := 0; ii < 500; ii++ {
					msg := ws.Message{
						MsgType:    "TEXT_MSG",
						SenderID:   "x",
						SenderName: "y",
						//Channels   []string `json:"channels"`
						Channel: "main",
						Content: fmt.Sprintf("MSG { clientID: %v, msgID: %v }", strID, ii),
					}
					connection.Receive(msg)

					r := rand.Intn(clients)
					if r%2 == 0 {
						// remo from chan
						gg := channels.RemoveClientFromChannel("channel-"+strconv.Itoa(r), client)
						if len(gg) > 0 {
							logger.Infof("--------", "-------", "-----------------------------", gg)
						}
					} else {
						// add
						channels.AddClientToChannel("channel-"+strconv.Itoa(r), client)
					}

				}
				msg3 := ws.Message{
					MsgType:    "LOGOUT_USER",
					SenderID:   "x",
					SenderName: "y",
					//Channels   []string `json:"channels"`
					Channel: "main",
					//Content: fmt.Sprintf("MSG { clientID: %v, msgID: %v }", strID, ii),
				}
				connection.Receive(msg3)
				logger.Infof("Messages", "End", "%v - finished sending messages", client)
			}(connection, strID, channels)

			time.Sleep(time.Second * 20)

			logger.Infof("Client", "Stop", "%v - finished sending messages", client)
			client.Stop()
			time.Sleep(time.Second * 1)
		}(i, channels)

	}
	time.Sleep(time.Second * 30)
}

func NewMockConnection() *MockConnection {
	c := make(chan ws.CommunicationResult, 10)
	return &MockConnection{
		communication: c,
	}
}

type MockConnection struct {
	communication chan ws.CommunicationResult
	channels      ws.Channels
}

func (c *MockConnection) Receive(msg ws.Message) {
	c.communication <- ws.CommunicationResult{
		Message: msg,
		Error:   nil,
	}
}

func (c *MockConnection) Send(msg ws.Message) error {
	//logger.Infof("MockConnection", "Send", "[Message] %v", msg)
	return nil
}

func (c *MockConnection) Close() error {
	logger.Info("MockConnection", "Close", "Close")
	return nil
}

func (c *MockConnection) Start() {

}

func (c *MockConnection) Incomming() chan ws.CommunicationResult {
	//logger.Info("MockConnection", "Incomming", "incomming")
	return c.communication
}
