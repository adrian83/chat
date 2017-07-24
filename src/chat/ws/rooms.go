package ws

import (
	"chat/logger"
	"fmt"
	"regexp"
)

var (
	roomNameRegexp = `^[a-zA-Z0-9_.-]*$`
	validRoomName  = regexp.MustCompile(roomNameRegexp)
)

// NewRooms returns new Rooms struct.
func NewRooms() Rooms {
	ch := make(map[string]Room)

	roomsListRequests := make(chan Client, 50)
	removeClient := make(chan Client, 50)
	addClientToRoomRequest := make(chan clientAndRoom, 50)
	removeClientFromRoomRequest := make(chan clientAndRoom, 50)
	createRoomRequest := make(chan clientAndRoom, 50)
	messageRequest := make(chan Message, 50)
	removeRoomRequests := make(chan string, 50)

	rooms := DefaultRooms{
		rooms:                       ch,
		roomsListRequests:           roomsListRequests,
		addClientToRoomRequest:      addClientToRoomRequest,
		removeClientFromRoomRequest: removeClientFromRoomRequest,
		createRoomRequest:           createRoomRequest,
		messageRequest:              messageRequest,
		removeRoomRequests:          removeRoomRequests,
		removeClient:                removeClient,
	}
	mainRoom := NewRoom(Main, &rooms)
	ch[Main] = mainRoom

	go rooms.start()

	return &rooms
}

// Rooms is an interface for describing collection of rooms.
type Rooms interface {
	ClientsRooms(client Client)
	RemoveRoom(name string)
	RemoveClient(client Client)
	CreateRoom(room string, client Client)
	AddClientToRoom(roomName string, client Client)
	RemoveClientFromRoom(roomName string, client Client)
	SendMessageOnRoom(message Message)
}

type clientAndRoom struct {
	client Client
	room   string
}

type messageAndRoom struct {
	message string
	room    string
}

// DefaultRooms struct represents collections of all rooms.
type DefaultRooms struct {
	rooms                       map[string]Room
	roomsListRequests           chan Client
	removeClient                chan Client
	removeRoomRequests          chan string
	addClientToRoomRequest      chan clientAndRoom
	removeClientFromRoomRequest chan clientAndRoom
	createRoomRequest           chan clientAndRoom
	messageRequest              chan Message
}

func (ch *DefaultRooms) start() {
	logger.Info("DefaultRooms", "start", "Starting DefaultRooms")

	for {
		select {
		case client := <-ch.roomsListRequests:

			rooms := make([]string, 0)
			for _, room := range ch.rooms {
				if cli := room.FindClient(client.ID()); cli != nil {
					rooms = append(rooms, room.Name())
				}
			}

			msg := RoomsNamesMessage(rooms)
			client.Send(msg)

		case cac := <-ch.addClientToRoomRequest:
			if room, exists := ch.rooms[cac.room]; exists {
				room.AddClient(cac.client)

				if cac.room == Main {
					names := make([]string, 0)
					for name := range ch.rooms {
						names = append(names, name)
					}
					roomNamesMsg := RoomsNamesMessage(names)
					cac.client.Send(roomNamesMsg)
				}

				ujc := NewUserJoinedRoomMessage(cac.room, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("DefaultRooms", "start", "Room %v does not exist. Client %v cannot join", cac.room, cac.client)
				errMsg := ErrorMessage(fmt.Sprintf("Room '%v' does not exist. Cannot join", cac.room))
				cac.client.Send(errMsg)
			}

		case roomName := <-ch.removeRoomRequests:
			if roomName != Main {
				delete(ch.rooms, roomName)
				msg := NewRemoveRoomMessage(roomName)
				mainRoom, _ := ch.rooms[Main]
				mainRoom.SendToEveryone(msg)
			} else {
				logger.Info("DefaultRooms", "start", "Cannot remove 'main' room")
			}

		case cac := <-ch.removeClientFromRoomRequest:
			logger.Infof("DefaultRooms", "start", "Remove client '%v' from room '%v'", cac.client, cac.room)
			if room, exists := ch.rooms[cac.room]; exists {
				room.RemoveClient(cac.client)
				ujc := NewUserLeftRoomMessage(cac.room, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("DefaultRooms", "start", "Room %v does not exist. Client %v cannot leave", cac.room, cac.client)
				errMsg := ErrorMessage(fmt.Sprintf("Room '%v' does not exist. Cannot leave", cac.room))
				cac.client.Send(errMsg)
			}

		case cac := <-ch.createRoomRequest:
			logger.Infof("DefaultRooms", "start", "Create room request from %v. Room name: %v", cac.client, cac.room)

			if cac.room == "" {
				errMsg := ErrorMessage("Invalid room name. Room name cannot be empty")
				cac.client.Send(errMsg)
				continue
			}

			if !validRoomName.MatchString(cac.room) {
				errMsg := ErrorMessage(fmt.Sprintf("Invalid room name. Name must match %v", roomNameRegexp))
				cac.client.Send(errMsg)
				continue
			}

			if _, exists := ch.rooms[cac.room]; !exists {
				// create new room with given name
				newRoom := NewRoom(cac.room, ch)
				newRoom.AddClient(cac.client)
				// add room to rooms' collection
				ch.rooms[cac.room] = newRoom

				ncm := NewAddRoomMessage(cac.room)

				mainRoom, _ := ch.rooms[Main]
				mainRoom.SendToEveryone(ncm)

				ujc := NewUserJoinedRoomMessage(cac.room, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("DefaultRooms", "start", "Room %v already exists. Client %v cannot create it", cac.room, cac.client)
			}

		case client := <-ch.removeClient:
			for _, room := range ch.rooms {
				room.RemoveClient(client)
			}

		case msg := <-ch.messageRequest:
			logger.Infof("DefaultRooms", "start", "Send message: %v", msg)

			if cha, exists := ch.rooms[msg.Room]; exists {
				cha.SendToEveryone(msg)
			} else {
				logger.Infof("DefaultRooms", "start", "Cannot send message because the room %v doesn't exist", msg.Room)
			}

		}
	}
}

// CreateRoom creates new request for creating new room.
func (ch *DefaultRooms) CreateRoom(roomName string, client Client) {
	ch.createRoomRequest <- clientAndRoom{
		client: client,
		room:   roomName,
	}
}

// RemoveClient removes client from all rooms.
func (ch *DefaultRooms) RemoveClient(client Client) {
	logger.Infof("DefaultRooms", "RemoveClient", "Removing Client %v from all rooms", client)
	ch.removeClient <- client
}

// RemoveRoom removes room with given name.
func (ch *DefaultRooms) RemoveRoom(roomName string) {
	ch.removeRoomRequests <- roomName
}

// ClientsRooms will return list of rooms to given client.
func (ch *DefaultRooms) ClientsRooms(client Client) {
	ch.roomsListRequests <- client
}

// AddClientToRoom adds given client to room with given name.
func (ch *DefaultRooms) AddClientToRoom(roomName string, client Client) {
	ch.addClientToRoomRequest <- clientAndRoom{
		client: client,
		room:   roomName,
	}
}

// RemoveClientFromRoom removes given client from room with given name.
func (ch *DefaultRooms) RemoveClientFromRoom(roomName string, client Client) {
	ch.removeClientFromRoomRequest <- clientAndRoom{
		client: client,
		room:   roomName,
	}
}

// SendMessageOnRoom sends given message to all clients of given room.
func (ch *DefaultRooms) SendMessageOnRoom(message Message) {
	ch.messageRequest <- message
}
