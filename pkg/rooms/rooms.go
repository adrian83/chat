package rooms

import (
	"fmt"
	"regexp"

	"github.com/adrian83/chat/pkg/exchange"

	logger "github.com/sirupsen/logrus"
)

var (
	roomNameRegexp = `^[a-zA-Z0-9_.-]*$`
	validRoomName  = regexp.MustCompile(roomNameRegexp)
)

// NewRooms returns new Rooms struct.
func NewRooms() *DefaultRooms {
	ch := make(map[string]*exchange.Room)

	roomsListRequests := make(chan exchange.Sender, 50)
	removeClient := make(chan exchange.Sender, 50)
	addClientToRoomRequest := make(chan clientAndRoom, 50)
	removeClientFromRoomRequest := make(chan clientAndRoom, 50)
	createRoomRequest := make(chan clientAndRoom, 50)
	messageRequest := make(chan exchange.Message, 50)
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
	mainRoom := exchange.NewMainRoom(&rooms)
	mainRoom.Start()
	ch[mainRoom.Name()] = mainRoom

	go rooms.start()

	return &rooms
}

type clientAndRoom struct {
	client exchange.Sender
	room   string
}

type messageAndRoom struct {
	message string
	room    string
}

type RoomsMap map[string]*exchange.Room

func (r RoomsMap) names() []string {
	names := make([]string, 0)
	for name := range r {
		names = append(names, name)
	}
	return names
}

// DefaultRooms struct represents collections of all rooms.
type DefaultRooms struct {
	rooms                       RoomsMap
	roomsListRequests           chan exchange.Sender
	removeClient                chan exchange.Sender
	removeRoomRequests          chan string
	addClientToRoomRequest      chan clientAndRoom
	removeClientFromRoomRequest chan clientAndRoom
	createRoomRequest           chan clientAndRoom
	messageRequest              chan exchange.Message
}

func (ch *DefaultRooms) start() {
	logger.Info("Starting DefaultRooms")

	for {
		select {
		case client := <-ch.roomsListRequests:
			rooms := ch.clientRooms(client.ID())
			msg := exchange.RoomsNamesMessage(rooms)
			client.Send(msg)

		case cac := <-ch.addClientToRoomRequest:
			if roomS, exists := ch.rooms[cac.room]; exists {
				roomS.AddClient(cac.client)

				if cac.room == exchange.MainRoomName() {
					names := ch.rooms.names()
					roomNamesMsg := exchange.RoomsNamesMessage(names)
					cac.client.Send(roomNamesMsg)
				}

				ujc := exchange.NewUserJoinedRoomMessage(cac.room, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("Room %v does not exist. Client %v cannot join", cac.room, cac.client)
				errMsg := exchange.ErrorMessage(fmt.Sprintf("Room '%v' does not exist. Cannot join", cac.room))
				cac.client.Send(errMsg)
			}

		case roomName := <-ch.removeRoomRequests:

			if roomName == exchange.MainRoomName() {
				logger.Info("Cannot remove 'main' room")
				continue
			}

			delete(ch.rooms, roomName)
			msg := exchange.NewRemoveRoomMessage(roomName)
			ch.sendToEveryone(exchange.MainRoomName(), msg)

		case cac := <-ch.removeClientFromRoomRequest:
			logger.Infof("Remove client '%v' from room '%v'", cac.client, cac.room)
			if room, exists := ch.rooms[cac.room]; exists {
				room.RemoveClient(cac.client.ID())
				ujc := exchange.NewUserLeftRoomMessage(cac.room, cac.client.ID())
				cac.client.Send(ujc)
			} else {
				logger.Infof("Room %v does not exist. Client %v cannot leave", cac.room, cac.client)
				errMsg := exchange.ErrorMessage(fmt.Sprintf("Room '%v' does not exist. Cannot leave", cac.room))
				cac.client.Send(errMsg)
			}

		case cac := <-ch.createRoomRequest:
			logger.Infof("Create room request from %v. Room name: %v", cac.client, cac.room)

			if cac.room == "" {
				errMsg := exchange.ErrorMessage("Invalid room name. Room name cannot be empty")
				cac.client.Send(errMsg)
				continue
			}

			if !validRoomName.MatchString(cac.room) {
				errMsg := exchange.ErrorMessage(fmt.Sprintf("Invalid room name. Name must match %v", roomNameRegexp))
				cac.client.Send(errMsg)
				continue
			}

			if _, exists := ch.rooms[cac.room]; exists {
				logger.Infof("Room %v already exists. Client %v cannot create it", cac.room, cac.client)
				continue
			}

			// create new room with given name
			newRoom := exchange.NewRoom(cac.room, ch)
			newRoom.Start()
			newRoom.AddClient(cac.client)
			// add room to rooms' collection
			ch.rooms[cac.room] = newRoom

			ncm := exchange.NewCreateRoomMessage(cac.room)
			ch.sendToEveryone(exchange.MainRoomName(), ncm)

			ujc := exchange.NewUserJoinedRoomMessage(cac.room, cac.client.ID())
			cac.client.Send(ujc)

		case client := <-ch.removeClient:
			for _, room := range ch.rooms {
				room.RemoveClient(client.ID())
			}

		case msg := <-ch.messageRequest:
			logger.Infof("Send message: %v", msg)
			ch.sendToEveryone(msg.Room, msg)
		}
	}
}

func (ch *DefaultRooms) sendToEveryone(roomName string, msg exchange.Message) {
	if room, ok := ch.rooms[msg.Room]; ok {
		room.SendToEveryone(msg)
	} else {
		logger.Infof("Cannot send message because the room %v doesn't exist", msg.Room)
	}
}

func (ch *DefaultRooms) clientRooms(id string) []string {
	rooms := make([]string, 0)
	for _, room := range ch.rooms {
		_, err := room.FindClient(id)
		if err == nil {
			rooms = append(rooms, room.Name())
		}
	}
	return rooms
}

// CreateRoom creates new request for creating new room.
func (ch *DefaultRooms) CreateRoom(roomName string, client exchange.Sender) {
	ch.createRoomRequest <- clientAndRoom{
		client: client,
		room:   roomName,
	}
}

// RemoveClient removes client from all rooms.
func (ch *DefaultRooms) RemoveClient(client exchange.Sender) {
	logger.Infof("Removing Client %v from all rooms", client)
	ch.removeClient <- client
}

// RemoveRoom removes room with given name.
func (ch *DefaultRooms) RemoveRoom(roomName string) {
	ch.removeRoomRequests <- roomName
}

// ClientsRooms will return list of rooms to given client.
func (ch *DefaultRooms) ClientsRooms(client exchange.Sender) {
	ch.roomsListRequests <- client
}

// AddClientToRoom adds given client to room with given name.
func (ch *DefaultRooms) AddClientToRoom(roomName string, client exchange.Sender) {
	ch.addClientToRoomRequest <- clientAndRoom{
		client: client,
		room:   roomName,
	}
}

// RemoveClientFromRoom removes given client from room with given name.
func (ch *DefaultRooms) RemoveClientFromRoom(roomName string, client exchange.Sender) {
	ch.removeClientFromRoomRequest <- clientAndRoom{
		client: client,
		room:   roomName,
	}
}

// SendMessageOnRoom sends given message to all clients of given room.
func (ch *DefaultRooms) SendMessageOnRoom(message exchange.Message) {
	ch.messageRequest <- message
}
