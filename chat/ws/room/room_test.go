package room

import (
	"testing"
	"time"

	"github.com/adrian83/chat/chat/ws/message"

	"github.com/stretchr/testify/assert"
)

func TestShouldCompareMainRoomName(t *testing.T) {
	assert.Equal(t, main, MainRoomName(), "Invalid main room name")
}

func TestShouldCreateNewRoomWithGivenName(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}
	name := "golang"

	room := NewRoom(name, rooms)

	assert.Equal(t, name, room.Name(), "Invalid room name")
	assert.False(t, room.Main(), "Room should not be main")
}

func TestShouldCreateNewMainRoom(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}

	room := NewMainRoom(rooms)

	assert.Equal(t, MainRoomName(), room.Name(), "Invalid main room name")
	assert.True(t, room.Main(), "Room should be main")
}

func TestCallingFindClientOnNotStartedRoomShouldSuspendExecution(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}
	name := "golang"

	room := NewRoom(name, rooms)

	finished := make(chan bool, 5)
	go func(room *Room, finished chan bool) {
		room.FindClient("test id")
		finished <- true
	}(room, finished)

	successAfterMs(t, 200, finished)
}

func TestShouldReturnErrorIfSearchedClientDoesNotExist(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}
	name := "golang"

	room := NewRoom(name, rooms)
	room.Start()

	client, err := room.FindClient("test id")

	assert.Empty(t, client, "Client should be nil")
	assert.Error(t, err, "Error should not be nil")
}

func TestShouldAddAndFindClient(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}
	sender := &SenderStub{id: "abc-def"}
	name := "golang"

	room := NewRoom(name, rooms)

	room.Start()
	time.Sleep(time.Duration(50) * time.Millisecond)

	room.AddClient(sender)
	time.Sleep(time.Duration(50) * time.Millisecond)

	client, err := room.FindClient(sender.id)

	assert.Equal(t, sender.ID(), client.ID(), "Invalid client id")
	assert.Empty(t, err, "Error should be nil")
}

type RoomsStub struct {
	removedRooms []string
}

func (r *RoomsStub) RemoveRoom(name string) {
	r.removedRooms = append(r.removedRooms, name)
}

type SenderStub struct {
	messages []message.Message
	id       string
	idCalled int
}

func (s *SenderStub) Send(msg message.Message) {
	s.messages = append(s.messages, msg)
}

func (s *SenderStub) ID() string {
	s.idCalled++
	return s.id
}

func successAfterMs(t *testing.T, ms int64, finished chan bool) {
	select {
	case <-finished:
		t.Error("Invalid state")
	case <-time.After(time.Duration(ms) * time.Millisecond):
		t.Log("success")
	}
}
