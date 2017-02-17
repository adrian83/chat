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

type Channels struct {
	channels map[string]*Channel
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
