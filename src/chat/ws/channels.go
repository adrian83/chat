package ws

import (
	"chat/logger"
	"fmt"
	"regexp"
)

var (
	channelNameRegexp = `^[a-zA-Z0-9_.-]*$`
	validChannelName  = regexp.MustCompile(channelNameRegexp)
)

// NewChannels returns new Channels struct.
func NewChannels() Channels {
	ch := make(map[string]Channel)

	interrupt := make(chan bool, 5)
	channelsListRequests := make(chan Client, 50)
	addClientToChannelRequest := make(chan clientAndChannel, 50)
	removeClientFromChannelRequest := make(chan clientAndChannel, 50)
	createChannelRequest := make(chan clientAndChannel, 50)
	messageRequest := make(chan Message, 50)

	channels := DefaultChannels{
		channels:                       ch,
		channelsListRequests:           channelsListRequests,
		interrupt:                      interrupt,
		addClientToChannelRequest:      addClientToChannelRequest,
		removeClientFromChannelRequest: removeClientFromChannelRequest,
		createChannelRequest:           createChannelRequest,
		messageRequest:                 messageRequest,
	}
	mainChannel := NewMainChannel(&channels)
	ch[Main] = mainChannel

	go channels.start()

	return &channels
}

// Channels is an interface for describing collection of channels.
type Channels interface {
	ClientsChannels(client Client)
	RemoveChannel(name string)
	RemoveClient(client Client)
	CreateChannel(channel string, client Client)
	AddClientToChannel(channelName string, client Client)
	RemoveClientFromChannel(channelName string, client Client)
	SendMessageOnChannel(message Message)
}

type clientAndChannel struct {
	client  Client
	channel string
}

type messageAndChannel struct {
	message string
	channel string
}

// DefaultChannels struct represents collections of all channels.
type DefaultChannels struct {
	channels                       map[string]Channel
	interrupt                      chan bool
	channelsListRequests           chan Client
	removeChannelRequests          chan string
	addClientToChannelRequest      chan clientAndChannel
	removeClientFromChannelRequest chan clientAndChannel
	createChannelRequest           chan clientAndChannel
	messageRequest                 chan Message
}

func (ch *DefaultChannels) start() {
	logger.Info("DefaultChannels", "start", "Starting DefaultChannels")

mainLoop:
	for {
		select {
		case client := <-ch.channelsListRequests:

			channels := make([]string, 0)
			for _, channel := range ch.channels {
				if _, exists := channel.FindClient(client.ID()); exists {
					channels = append(channels, channel.Name())
				}
			}

			msg := ChannelsNamesMessage(channels)
			client.Send(msg)

		case cac := <-ch.addClientToChannelRequest:
			if channel, exists := ch.channels[cac.channel]; exists {
				channel.AddClient(cac.client)

				if cac.channel == Main {
					names := make([]string, 0)
					for name := range ch.channels {
						names = append(names, name)
					}
					channelNamesMsg := ChannelsNamesMessage(names)
					cac.client.Send(channelNamesMsg)
				}

				ujc := NewUserJoinedChannelMessage(cac.channel, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("DefaultChannels", "start", "Channel %v does not exist. Client %v cannot join", cac.channel, cac.client)
				errMsg := ErrorMessage(fmt.Sprintf("Channel '%v' does not exist. Cannot join", cac.channel))
				cac.client.Send(errMsg)
			}

		case channelName := <-ch.removeChannelRequests:
			if channelName != Main {
				delete(ch.channels, channelName)
				msg := NewRemoveChannelMessage(channelName)
				mainRoom, _ := ch.channels[Main]
				mainRoom.SendToEveryone(msg)
			} else {
				logger.Info("DefaultChannels", "start", "Cannot remove 'main' room")
			}

		case cac := <-ch.removeClientFromChannelRequest:
			logger.Infof("DefaultChannels", "start", "Remove client '%v' from room '%v'", cac.client, cac.channel)
			if channel, exists := ch.channels[cac.channel]; exists {
				channel.RemoveClient(cac.client)
				ujc := NewUserLeftChannelMessage(cac.channel, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("DefaultChannels", "start", "Channel %v does not exist. Client %v cannot leave", cac.channel, cac.client)
				errMsg := ErrorMessage(fmt.Sprintf("Channel '%v' does not exist. Cannot leave", cac.channel))
				cac.client.Send(errMsg)
			}

		case cac := <-ch.createChannelRequest:
			logger.Infof("DefaultChannels", "start", "Create channel request from %v. Channel name: %v", cac.client, cac.channel)

			if !validChannelName.MatchString(cac.channel) {
				errMsg := ErrorMessage(fmt.Sprintf("Invalid channel name. Name must match %v", channelNameRegexp))
				cac.client.Send(errMsg)
				continue
			}

			if _, exists := ch.channels[cac.channel]; !exists {
				// create new room with given name
				newRoom := NewChannel(cac.channel, cac.client, ch)
				// add room to rooms' collection
				ch.channels[cac.channel] = newRoom

				ncm := NewAddChannelMessage(cac.channel)

				mainRoom, _ := ch.channels[Main]
				mainRoom.SendToEveryone(ncm)

				ujc := NewUserJoinedChannelMessage(cac.channel, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("DefaultChannels", "start", "Channel %v already exists. Client %v cannot create it", cac.channel, cac.client)
			}

		case msg := <-ch.messageRequest:
			logger.Infof("DefaultChannels", "start", "Send message: %v", msg)

			if cha, exists := ch.channels[msg.Channel]; exists {
				cha.SendToEveryone(msg)
			} else {
				logger.Infof("DefaultChannels", "start", "Cannot send message because the channel %v doesn't exist", msg.Channel)
			}

		case <-ch.interrupt:
			break mainLoop

		}
	}

	logger.Info("DefaultChannels", "start", "Stopping DefaultChannels")
}

// CreateChannel creates new request for creating new channel.
func (ch *DefaultChannels) CreateChannel(channelName string, client Client) {
	ch.createChannelRequest <- clientAndChannel{
		client:  client,
		channel: channelName,
	}
}

// RemoveClient removes client from all channels.
func (ch *DefaultChannels) RemoveClient(client Client) {
	logger.Infof("DefaultChannels", "RemoveClient", "Removing Client %v from all channels (NOT IMPLEMENTED)", client)
}

// RemoveChannel removes room with given name.
func (ch *DefaultChannels) RemoveChannel(channelName string) {
	ch.removeChannelRequests <- channelName
}

// ClientsChannels will return list of channels to given client.
func (ch *DefaultChannels) ClientsChannels(client Client) {
	ch.channelsListRequests <- client
}

// AddClientToChannel adds given client to channel with given name.
func (ch *DefaultChannels) AddClientToChannel(channelName string, client Client) {
	ch.addClientToChannelRequest <- clientAndChannel{
		client:  client,
		channel: channelName,
	}
}

// RemoveClientFromChannel removes given client from channel with given name.
func (ch *DefaultChannels) RemoveClientFromChannel(channelName string, client Client) {
	ch.removeClientFromChannelRequest <- clientAndChannel{
		client:  client,
		channel: channelName,
	}
}

// SendMessageOnChannel sends given message to all clients of given channel.
func (ch *DefaultChannels) SendMessageOnChannel(message Message) {
	ch.messageRequest <- message
}
