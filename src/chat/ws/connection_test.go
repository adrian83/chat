package ws

// MockConnection implementation of Connection interface used for testing.
type MockConnection struct {
	LastMesssage Message
	Closed       bool
	//Send(msg Message) error
	//Close() error
	//Start()
	//Incomming() chan CommunicationResult
}

func NewMockConnection() *MockConnection {
	return &MockConnection{
		LastMesssage: Message{},
		Closed:       false,
	}
}

func (c *MockConnection) Send(msg Message) error {
	c.LastMesssage = msg
	return nil
}

func (c *MockConnection) Close() error {
	c.Closed = true
	return nil
}

func (c *MockConnection) Start() {

}
func (c *MockConnection) Incomming() chan CommunicationResult {
	return make(chan CommunicationResult)
}
