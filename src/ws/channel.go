package ws

const (
	main = "main"
)

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

type Channel struct {
	name    string
	clients map[string]*Client
}

func NewChannel(name string, client *Client) *Channel {
	return &Channel{
		name:    name,
		clients: map[string]*Client{client.id: client},
	}
}

type Channels struct {
	channels map[string]*Channel
}

func (ch *Channels) Names() []string {
	names := make([]string, 0)
	for name, _ := range ch.channels {
		names = append(names, name)
	}
	return names
}

func (ch *Channels) AddChannel(channel *Channel) {
	ch.channels[channel.name] = channel
}

func (ch *Channels) RegisterClient(client *Client) {
	ch.channels[main].clients[client.id] = client
	client.AddChannel(ch.channels[main])
}

func (ch *Channels) RemoveClient(client *Client) {
	for _, channel := range ch.channels {
		delete(channel.clients, client.id)
	}
}
