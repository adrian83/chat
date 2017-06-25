package ws

import (
	//"logger"
	"sync"
)

// NewChannels returns new Channels struct.
func NewChannels() Channels {
	ch := make(map[string]Channel)
	ch[main] = &DefaultChannel{
		name:    main,
		clients: make(map[string]Client),
	}
	channels := DefaultChannels{
		channels: ch,
	}
	return &channels
}

// Channels is an interface for describing collection of channels.
type Channels interface {
	Names() []string
	GetMainChannel() Channel
	ClientsChannels(client Client) map[string]Channel
	AddChannel(channel Channel) []SendError
	RemoveChannel(name string)
	FindChannel(name string) (Channel, bool)
	AddClientToChannel(channelName string, client Client)
	RemoveClientFromChannel(channelName string, client Client) []SendError
	AddClient(client Client) error
	RemoveClient(client Client)
}

// DefaultChannels struct represents collections of all channels.
type DefaultChannels struct {
	lock     sync.RWMutex
	channels map[string]Channel
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

// ClientsChannels returns all clients channels.
func (ch *DefaultChannels) ClientsChannels(client Client) map[string]Channel {
	channels := make(map[string]Channel, 0)
	ch.lock.RLock()
	for _, channel := range ch.channels {
		if _, exists := channel.FindClient(client.ID()); exists {
			channels[channel.Name()] = channel
		}
	}
	ch.lock.RUnlock()
	return channels
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
			if err := client.Send(msg); err != nil {
				sendErrors = append(sendErrors, SendError{
					Client: client,
					Err:    err,
				})
				continue
			}
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
	if channel, ok := ch.FindChannel(channelName); ok {
		channel.AddClient(client)
	}
}

// RemoveClientFromChannel removes given client from channel with given name.
func (ch *DefaultChannels) RemoveClientFromChannel(channelName string, client Client) []SendError {
	if channelName != main {
		if channel, ok := ch.FindChannel(channelName); ok {
			return channel.RemoveClient(client)
		}
	}
	return make([]SendError, 0)
}

// AddClient registers new client and sends him some information
func (ch *DefaultChannels) AddClient(client Client) error {
	ch.lock.RLock()
	ch.channels[main].AddClient(client)
	ch.lock.RUnlock()

	channelNamesMsg := ChannelsNamesMessage(ch.Names())
	return client.Send(channelNamesMsg)
}

// RemoveClient removes client.
func (ch *DefaultChannels) RemoveClient(client Client) {
	ch.lock.RLock()
	for _, channel := range ch.channels {
		channel.RemoveClientFromChannel(client)
	}
	ch.lock.RUnlock()
}
