package ws

import (
	"chat/logger"
	"sync"
)

// NewChannels returns new Channels struct.
func NewChannels() Channels {
	ch := make(map[string]Channel)

	interrupt := make(chan bool, 5)
	channelsListRequests := make(chan Client, 50)
	addClientToChannelRequest := make(chan clientAndChannel, 50)
	createChannelRequest := make(chan clientAndChannel, 50)
	messageRequest := make(chan Message, 50)

	channels := DefaultChannels{
		channels:                  ch,
		channelsListRequests:      channelsListRequests,
		interrupt:                 interrupt,
		addClientToChannelRequest: addClientToChannelRequest,
		createChannelRequest:      createChannelRequest,
		messageRequest:            messageRequest,
	}
	mainChannel := NewMainChannel(&channels)
	ch[main] = mainChannel

	go channels.start()

	return &channels
}

// Channels is an interface for describing collection of channels.
type Channels interface {
	Names() []string
	GetMainChannel() Channel
	ClientsChannels(client Client)
	AddChannel(channel Channel) []SendError
	RemoveChannel(name string)
	FindChannel(name string) (Channel, bool)

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
	lock     sync.RWMutex
	channels map[string]Channel

	interrupt                      chan bool
	channelsListRequests           chan Client
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
			if channel, exists := ch.FindChannel(cac.channel); exists {
				channel.AddClient(cac.client)
				ujc := NewUserJoinedChannelMessage(cac.channel, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("DefaultChannels", "start", "Channel %v does not exist. Client %v cannot join", cac.channel, cac.client)
			}

		case cac := <-ch.removeClientFromChannelRequest:
			if channel, exists := ch.FindChannel(cac.channel); exists {
				channel.RemoveClient(cac.client)
				ujc := NewUserLeftChannelMessage(cac.channel, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("DefaultChannels", "start", "Channel %v does not exist. Client %v cannot leave", cac.channel, cac.client)
			}

		case cac := <-ch.createChannelRequest:
			logger.Infof("DefaultChannels", "start", "Create channel request from %v. Channel name: %v", cac.client, cac.channel)

			if _, exists := ch.FindChannel(cac.channel); !exists {
				cha := NewChannel(cac.channel, cac.client, ch)
				ch.AddChannel(cha)
				//ncm := NewAddChannelMessage(cac.channel)
				//cac.client.Send(ncm)
			} else {
				logger.Infof("DefaultChannels", "start", "Channel %v already exists. Client %v cannot create it", cac.channel, cac.client)
			}

		case msg := <-ch.messageRequest:
			logger.Infof("DefaultChannels", "start", "Send message: %v", msg)

			if cha, exists := ch.FindChannel(msg.Channel); exists {
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
	logger.Info("DefaultChannels", "CreateChannel", "CreateChannel")
	ch.createChannelRequest <- clientAndChannel{
		client:  client,
		channel: channelName,
	}
}

// RemoveClient removes client from all channels.
func (ch *DefaultChannels) RemoveClient(client Client) {
	logger.Infof("DefaultChannels", "RemoveClient", "Removing Client %v from all channels (NOT IMPLEMENTED)", client)
}

// GetMainChannel returns main channel.
func (ch *DefaultChannels) GetMainChannel() Channel {
	ch.lock.RLock()
	mainChannel := ch.channels[main]
	ch.lock.RUnlock()
	return mainChannel
}

// RemoveChannel removes channel.
func (ch *DefaultChannels) RemoveChannel(name string) {
	if name != main {
		ch.lock.Lock()
		delete(ch.channels, name)
		ch.lock.Unlock()
	}
}

// ClientsChannels will return list of channels to given client.
func (ch *DefaultChannels) ClientsChannels(client Client) {
	ch.channelsListRequests <- client
}

// Names return names of all channels.
func (ch *DefaultChannels) Names() []string {
	names := make([]string, 0)
	ch.lock.RLock()
	for name := range ch.channels {
		names = append(names, name)
	}
	ch.lock.RUnlock()
	return names
}

// AddChannel add channel to collection of channels.
func (ch *DefaultChannels) AddChannel(channel Channel) []SendError {

	msg := NewAddChannelMessage(channel.Name())

	// check if channel with this name already exists
	existingChannel, exists := ch.FindChannel(channel.Name())

	if exists {
		// if channel exists, add every client from given channel to existing one
		sendErrors := make([]SendError, 0)
		for name, client := range channel.Clients() {
			client.Send(msg)
			existingChannel.Clients()[name] = client
		}
		return sendErrors
	}

	// if channel with such name doesn't exist, add this channel
	ch.lock.Lock()
	ch.channels[channel.Name()] = channel
	ch.lock.Unlock()

	mainChannel := ch.GetMainChannel()

	// there is only one client in the channel and he should be
	// set as a sender of the message
	for _, c := range channel.Clients() {
		msg.SenderID = c.ID()
		break
	}
	return mainChannel.SendToEveryone(msg)

}

// FindChannel returns channel with given name it it exist.
func (ch *DefaultChannels) FindChannel(channelName string) (Channel, bool) {
	ch.lock.RLock()
	channel, ok := ch.channels[channelName]
	ch.lock.RUnlock()
	return channel, ok
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
