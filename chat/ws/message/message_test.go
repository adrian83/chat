package message

import "testing"

func TestUnknownMessageImplementsHandler(t *testing.T) {
	var _ Handler = &UnknownMessage{}
}

type SenderStub struct {
	id string
}

func (s *SenderStub) Send(msg Message) {

}
func (s *SenderStub) ID() string {
	return s.id
}

type RoomsStub struct {
}

func (r *RoomsStub) AddClientToRoom(string, Sender)      {}
func (r *RoomsStub) RemoveClientFromRoom(string, Sender) {}
func (r *RoomsStub) SendMessageOnRoom(Message)           {}
func (r *RoomsStub) CreateRoom(string, Sender)           {}
