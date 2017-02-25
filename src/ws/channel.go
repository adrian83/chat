package ws

import (
	"logger"
)

const (
	main = "main"
)

// SendError error representing error while sending msg to client.
type SendError struct {
	Client Client
	Err    error
}

// Error implementation of error interface.
func (e SendError) Error() string {
	return e.Error()
}

// Channel is an interface for defining channels.
type Channel interface {
	Name() string
	Clients() map[string]Client
	SendToEveryone(msg Message) []SendError
	RemoveClient(client Client) []SendError
}

// DefaultChannel struct representing chat channel.
type DefaultChannel struct {
	name     string
	clients  map[string]Client
	channels *Channels
}

// Empty returns true if there is no clients in this channel, false otherwise
func (ch *DefaultChannel) Empty() bool {
	return 0 == len(ch.clients)
}

// Name returns chanel's name.
func (ch *DefaultChannel) Name() string {
	return ch.name
}

// Clients returns chanel's clients.
func (ch *DefaultChannel) Clients() map[string]Client {
	return ch.clients
}

// SendToEveryone sends message to everyone from that channel.
func (ch *DefaultChannel) SendToEveryone(msg Message) []SendError {
	logger.Infof("Channel", "SendToEveryone", "Sending msg to everyone from channel '%v'.", ch.name)
	errors := make([]SendError, 0)
	for _, client := range ch.clients {
		logger.Infof("Channel", "SendToEveryone", "Sending msg to %v from channel '%v'.", client, ch.name)
		if err := client.Send(msg); err != nil {
			sendErr := SendError{
				Client: client,
				Err:    err,
			}
			errors = append(errors, sendErr)
		}
	}
	return errors
}

// RemoveClient removes client from the channel.
func (ch *DefaultChannel) RemoveClient(client Client) []SendError {
	if ch.name != main {
		delete(ch.clients, client.ID())
		if ch.Empty() {
			msg := NewRemoveChannelMessage(ch.name)
			mainChannel := ch.channels.channels[main]
			return mainChannel.SendToEveryone(msg)
		}
	}
	return make([]SendError, 0)
}

// NewChannel functions returns new Channel struct.
func NewChannel(name string, client Client, channels *Channels) Channel {
	return &DefaultChannel{
		name:     name,
		clients:  map[string]Client{client.ID(): client},
		channels: channels,
	}
}

// NewChannels returns new Channels struct.
func NewChannels() *Channels {
	ch := make(map[string]Channel)
	ch[main] = &DefaultChannel{
		name:    main,
		clients: make(map[string]Client),
	}
	channels := Channels{
		channels: ch,
	}
	return &channels
}

// Channels struct represents collections of ll channels.
type Channels struct {
	channels map[string]Channel
}

// ClientsChannels returns all clients channels .
func (ch *Channels) ClientsChannels(client Client) map[string]Channel {
	chs := make(map[string]Channel, 0)
	for _, ch := range ch.channels {
		if _, ok := ch.Clients()[client.ID()]; ok {
			chs[ch.Name()] = ch
		}
	}
	return chs
}

// Names return names of all channels.
func (ch *Channels) Names() []string {
	names := make([]string, 0)
	for name := range ch.channels {
		names = append(names, name)
	}
	return names
}

// AddChannel add channel to collection of channels.
func (ch *Channels) AddChannel(channel Channel) []SendError {
	errs := make([]SendError, 0)
	msg := NewAddChannelMessage(channel.Name())
	if cha, ok := ch.channels[channel.Name()]; ok {
		for name, c := range channel.Clients() {
			cha.Clients()[name] = c
			sendErr := SendError{
				Client: c,
				Err:    c.Send(msg),
			}
			errs = append(errs, sendErr)
		}
	} else {
		ch.channels[channel.Name()] = channel
		mainChannel := ch.channels[main]
		for _, c := range channel.Clients() {
			msg.SenderID = c.ID()
		}
		return mainChannel.SendToEveryone(msg)
	}
	return errs
}

// AddClientToChannel adds given client to channel with given name.
func (ch *Channels) AddClientToChannel(channelName string, client Client) {
	ch.channels[channelName].Clients()[client.ID()] = client
}

// RemoveClientFromChannel removes given client from channel with given name.
func (ch *Channels) RemoveClientFromChannel(channelName string, client Client) []SendError {
	return ch.channels[channelName].RemoveClient(client)
}

// RegisterClient registers new client and sends him some information
func (ch *Channels) RegisterClient(client Client) error {
	ch.channels[main].Clients()[client.ID()] = client
	channelNamesMsg := ChannelsNamesMessage(ch.Names())
	return client.Send(channelNamesMsg)
}

// RemoveClient removes client.
func (ch *Channels) RemoveClient(client Client) {
	for _, channel := range ch.channels {
		delete(channel.Clients(), client.ID())
	}
}
