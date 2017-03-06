package ws

import (
	//"logger"
	"sync"
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
	FindClient(clientID string) (Client, bool)
	SendToEveryone(msg Message) []SendError
	RemoveClient(client Client) []SendError
	RemoveClientFromChannel(client Client)
	AddClient(client Client)
}

// DefaultChannel struct representing chat channel.
type DefaultChannel struct {
	name     string
	lock     sync.RWMutex
	clients  map[string]Client
	channels *Channels
}

// Empty returns true if there is no clients in this channel, false otherwise
func (ch *DefaultChannel) Empty() bool {
	ch.lock.RLock()
	result := 0 == len(ch.clients)
	ch.lock.RUnlock()
	return result
}

// FindClient returns client with given id if it exist in this channel.
func (ch *DefaultChannel) FindClient(clientID string) (Client, bool) {
	ch.lock.RLock()
	c, ok := ch.clients[clientID]
	ch.lock.RUnlock()
	return c, ok
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
	//logger.Infof("Channel", "SendToEveryone", "Sending msg to everyone from channel '%v'.", ch.name)
	errors := make([]SendError, 0)
	ch.lock.RLock()
	for _, client := range ch.clients {
		//logger.Infof("Channel", "SendToEveryone", "Sending msg to %v from channel '%v'.", client, ch.name)
		if err := client.Send(msg); err != nil {
			sendErr := SendError{
				Client: client,
				Err:    err,
			}
			errors = append(errors, sendErr)
		}
	}
	ch.lock.RUnlock()
	return errors
}

// RemoveClientFromChannel removes client from the channel.
func (ch *DefaultChannel) RemoveClientFromChannel(client Client) {
	if ch.name != main {
		ch.lock.Lock()
		delete(ch.clients, client.ID())
		ch.lock.Unlock()
	}
}

// AddClient adds client to channel.
func (ch *DefaultChannel) AddClient(client Client) {
	ch.lock.Lock()
	ch.clients[client.ID()] = client
	ch.lock.Unlock()
}

// RemoveClient removes client from the channel.
func (ch *DefaultChannel) RemoveClient(client Client) []SendError {
	if ch.name != main {

		ch.lock.Lock()
		delete(ch.clients, client.ID())
		ch.lock.Unlock()

		if ch.Empty() {
			msg := NewRemoveChannelMessage(ch.name)

			ch.channels.lock.RLock()
			mainChannel := ch.channels.channels[main]
			ch.channels.lock.RUnlock()

			ch.channels.lock.Lock()
			delete(ch.channels.channels, ch.Name())
			ch.channels.lock.Unlock()

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
	lock     sync.RWMutex
	channels map[string]Channel
}

// ClientsChannels returns all clients channels .
func (ch *Channels) ClientsChannels(client Client) map[string]Channel {
	chs := make(map[string]Channel, 0)
	ch.lock.RLock()
	for _, ch := range ch.channels {
		if _, ok := ch.FindClient(client.ID()); ok {
			chs[ch.Name()] = ch
		}
	}
	ch.lock.RUnlock()
	return chs
}

// Names return names of all channels.
func (ch *Channels) Names() []string {
	names := make([]string, 0)
	ch.lock.RLock()
	for name := range ch.channels {
		names = append(names, name)
	}
	ch.lock.RUnlock()
	return names
}

// AddChannel add channel to collection of channels.
func (ch *Channels) AddChannel(channel Channel) []SendError {
	errs := make([]SendError, 0)
	msg := NewAddChannelMessage(channel.Name())

	ch.lock.RLock()
	cha, ok := ch.channels[channel.Name()]
	ch.lock.RUnlock()

	if ok {
		for name, c := range channel.Clients() {
			cha.Clients()[name] = c
			sendErr := SendError{
				Client: c,
				Err:    c.Send(msg),
			}
			errs = append(errs, sendErr)
		}
	} else {
		ch.lock.Lock()
		ch.channels[channel.Name()] = channel
		ch.lock.Unlock()

		ch.lock.RLock()
		mainChannel := ch.channels[main]
		ch.lock.RUnlock()

		//channel.lock
		for _, c := range channel.Clients() {
			msg.SenderID = c.ID()
		}
		return mainChannel.SendToEveryone(msg)
	}
	return errs
}

// AddClientToChannel adds given client to channel with given name.
func (ch *Channels) AddClientToChannel(channelName string, client Client) {
	ch.lock.RLock()
	chann, ok := ch.channels[channelName]
	ch.lock.RUnlock()
	if ok {
		chann.AddClient(client)
	}
}

// RemoveClientFromChannel removes given client from channel with given name.
func (ch *Channels) RemoveClientFromChannel(channelName string, client Client) []SendError {
	if channelName != main {
		ch.lock.RLock()
		chann, ok := ch.channels[channelName]
		ch.lock.RUnlock()
		if ok {
			return chann.RemoveClient(client)
		}
	}
	return make([]SendError, 0)
}

// RegisterClient registers new client and sends him some information
func (ch *Channels) RegisterClient(client Client) error {
	ch.lock.RLock()
	ch.channels[main].AddClient(client)
	ch.lock.RUnlock()

	channelNamesMsg := ChannelsNamesMessage(ch.Names())
	return client.Send(channelNamesMsg)
}

// RemoveClient removes client.
func (ch *Channels) RemoveClient(client Client) {
	ch.lock.RLock()
	for _, channel := range ch.channels {
		channel.RemoveClientFromChannel(client)
	}
	ch.lock.RUnlock()
}
