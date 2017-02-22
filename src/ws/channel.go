package ws

import (
	"logger"
)

const (
	main = "main"
)

// SendError error representing error while sending msg to client.
type SendError struct {
	Client *Client
	Err    error
}

// Error implementation of error interface.
func (e SendError) Error() string {
	return e.Error()
}

// Channel struct representing chat channel.
type Channel struct {
	name     string
	clients  map[string]*Client
	channels *Channels
}

// Empty returns true if there is no clients in this channel, false otherwise
func (ch *Channel) Empty() bool {
	return 0 == len(ch.clients)
}

// SendToEveryone sends message to everyone from that channel.
func (ch *Channel) SendToEveryone(msg map[string]interface{}) []SendError {
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
func (ch *Channel) RemoveClient(client *Client) []SendError {
	if ch.name != main {
		delete(ch.clients, client.id)
		if ch.Empty() {
			msg := map[string]interface{}{"msgType": "REM_CH", "channel": ch.name, "senderId": "system"}
			mainChannel := ch.channels.channels[main]
			return mainChannel.SendToEveryone(msg)
		}
	}
	return make([]SendError, 0)
}

// NewChannel functions returns new Channel struct.
func NewChannel(name string, client *Client, channels *Channels) *Channel {
	return &Channel{
		name:     name,
		clients:  map[string]*Client{client.id: client},
		channels: channels,
	}
}

// NewChannels returns new Channels struct.
func NewChannels() *Channels {
	ch := make(map[string]*Channel)
	ch[main] = &Channel{
		name:    main,
		clients: make(map[string]*Client),
	}
	channels := Channels{
		channels: ch,
	}
	return &channels
}

// Channels struct represents collections of ll channels.
type Channels struct {
	channels map[string]*Channel
}

// ClientsChannels returns all clients channels .
func (ch *Channels) ClientsChannels(client *Client) map[string]*Channel {
	chs := make(map[string]*Channel, 0)
	for _, ch := range ch.channels {
		if _, ok := ch.clients[client.id]; ok {
			chs[ch.name] = ch
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
func (ch *Channels) AddChannel(channel *Channel) []SendError {
	errs := make([]SendError, 0)
	msg := map[string]interface{}{"msgType": "ADD_CH", "channel": channel.name, "senderId": "system"}
	if cha, ok := ch.channels[channel.name]; ok {
		for name, c := range channel.clients {
			cha.clients[name] = c
			sendErr := SendError{
				Client: c,
				Err:    c.Send(msg),
			}
			errs = append(errs, sendErr)
		}
	} else {
		ch.channels[channel.name] = channel
		mainChannel := ch.channels[main]
		for _, c := range channel.clients {
			msg["senderId"] = c.id
		}
		return mainChannel.SendToEveryone(msg)

	}
	return errs

}

// AddClientToChannel adds given client to channel with given name.
func (ch *Channels) AddClientToChannel(channelName string, client *Client) {
	ch.channels[channelName].clients[client.id] = client
}

func (ch *Channels) RegisterClient(client *Client) {
	ch.channels[main].clients[client.id] = client
}

func (ch *Channels) RemoveClient(client *Client) {
	for _, channel := range ch.channels {
		delete(channel.clients, client.id)
	}
}
