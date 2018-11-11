package room

import (
	"testing"

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
}

func TestShouldCreateNewMainRoom(t *testing.T) {
	rooms := &RoomsStub{removedRooms: make([]string, 0)}

	room := NewMainRoom(rooms)

	assert.Equal(t, MainRoomName(), room.Name(), "Invalid main room name")
}

type RoomsStub struct {
	removedRooms []string
}

func (r *RoomsStub) RemoveRoom(name string) {
	r.removedRooms = append(r.removedRooms, name)
}
