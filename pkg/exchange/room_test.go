package exchange

import (
	"testing"
	"time"

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

func TestSendingFewMessagesOnNonStartedRoom(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}

	room := NewMainRoom(rooms)

	finished := make(chan bool, 5)
	go func(room *Room, finished chan bool) {
		for i := 0; i < 20; i++ {
			room.SendToEveryone(NewCreateRoomMessage("java"))
		}
		finished <- true
	}(room, finished)

	failAfterMs(t, 200, finished)
}

func TestSendingMultipleMessagesOnNonStartedRoomShouldSuspendExecution(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}
	msg := NewCreateRoomMessage("java")

	room := NewMainRoom(rooms)

	finished := make(chan bool, 5)
	go func(room *Room, finished chan bool) {
		for i := 0; i < 60; i++ {
			room.SendToEveryone(msg)
		}
		finished <- true
	}(room, finished)

	successAfterMs(t, 200, finished)
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

func TestShouldRemoveClientAndRoom(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}
	sender := &SenderStub{id: "abc-def"}
	name := "golang"

	room := NewRoom(name, rooms)

	room.Start()
	time.Sleep(time.Duration(50) * time.Millisecond)

	room.AddClient(sender)
	time.Sleep(time.Duration(50) * time.Millisecond)

	room.RemoveClient(sender.ID())
	time.Sleep(time.Duration(50) * time.Millisecond)

	assert.Equal(t, 1, len(rooms.removedRooms), "One room should be removed")
}

func TestShouldRemoveClientButNotCloseRoom(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}
	sender1 := &SenderStub{id: "abc-def"}
	sender2 := &SenderStub{id: "ghi-jkl"}
	name := "golang"

	room := NewRoom(name, rooms)

	room.Start()
	time.Sleep(time.Duration(50) * time.Millisecond)

	room.AddClient(sender1)
	room.AddClient(sender2)
	time.Sleep(time.Duration(50) * time.Millisecond)

	room.RemoveClient(sender2.ID())
	time.Sleep(time.Duration(50) * time.Millisecond)

	client1, err := room.FindClient(sender1.id)

	assert.Equal(t, sender1.ID(), client1.ID(), "Invalid client id")
	assert.Empty(t, err, "Error should be nil")

	client2, err := room.FindClient(sender2.id)

	assert.Empty(t, client2, "Client should be nil")
	assert.Error(t, err, "Error should not be nil")

	assert.Equal(t, 0, len(rooms.removedRooms), "No rooms should be removed")
}

func TestShouldSendMultipleMessages(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}
	sender1 := &SenderStub{id: "abc-def"}
	sender2 := &SenderStub{id: "ghi-jkl"}
	msg := NewCreateRoomMessage("java")
	name := "golang"
	msgNumber := 60

	room := NewRoom(name, rooms)

	room.Start()
	time.Sleep(time.Duration(50) * time.Millisecond)

	room.AddClient(sender1)
	room.AddClient(sender2)
	time.Sleep(time.Duration(50) * time.Millisecond)

	for i := 0; i < msgNumber; i++ {
		room.SendToEveryone(msg)
	}
	time.Sleep(time.Duration(50) * time.Millisecond)

	assert.Equal(t, msgNumber, len(sender1.messages), "invalid number of received messages")
	assert.Equal(t, msgNumber, len(sender2.messages), "invalid number of received messages")
}

type SenderStub struct {
	messages     []Message
	id           string
	idCalled     int
	sendExecuted int
}

func (s *SenderStub) Send(msg Message) {
	s.sendExecuted++
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

func failAfterMs(t *testing.T, ms int64, finished chan bool) {
	select {
	case <-finished:
		t.Log("success")
	case <-time.After(time.Duration(ms) * time.Millisecond):
		t.Error("Invalid state")
	}
}
