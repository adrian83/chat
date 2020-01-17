package exchange

import (
	"regexp"

	logger "github.com/sirupsen/logrus"
)

var (
	roomNameRegexp = `^[a-zA-Z0-9_.-]*$`
	validRoomName  = regexp.MustCompile(roomNameRegexp)
)

// NewRooms returns new Rooms struct.
func NewRooms() *Rooms {
	ch := make(map[string]*Room)

	roomsListRequests := make(chan *Client, 50)
	removeClient := make(chan *Client, 50)
	addClientToRoomRequest := make(chan clientAndRoom, 50)
	removeClientFromRoomRequest := make(chan clientAndRoom, 50)
	createRoomRequest := make(chan clientAndRoom, 50)
	messageRequest := make(chan *Message, 50)
	removeRoomRequests := make(chan string, 50)

	rooms := Rooms{
		rooms:                       ch,
		roomsListRequests:           roomsListRequests,
		addClientToRoomRequest:      addClientToRoomRequest,
		removeClientFromRoomRequest: removeClientFromRoomRequest,
		createRoomRequest:           createRoomRequest,
		messageRequest:              messageRequest,
		removeRoomRequests:          removeRoomRequests,
		removeClient:                removeClient,
	}
	mainRoom := NewMainRoom(&rooms)
	mainRoom.Start()
	ch[mainRoom.Name()] = mainRoom

	go rooms.start()

	return &rooms
}

type clientAndRoom struct {
	client *Client
	room   string
}

type RoomsMap map[string]*Room

func (r RoomsMap) names() []string {
	names := make([]string, 0)
	for name := range r {
		names = append(names, name)
	}
	return names
}

// Rooms struct represents collections of all rooms.
type Rooms struct {
	rooms                       RoomsMap
	roomsListRequests           chan *Client
	removeClient                chan *Client
	removeRoomRequests          chan string
	addClientToRoomRequest      chan clientAndRoom
	removeClientFromRoomRequest chan clientAndRoom
	createRoomRequest           chan clientAndRoom
	messageRequest              chan *Message
}

func (ch *Rooms) start() {
	logger.Info("Starting Rooms")

	for {
		select {
		case client := <-ch.roomsListRequests:
			rooms := ch.clientRooms(client.ID())
			msg := RoomsNamesMessage(rooms)
			client.Send(msg)

		case cac := <-ch.addClientToRoomRequest:
			roomS := ch.rooms[cac.room]
			roomS.AddClient(cac.client)

			names := ch.rooms.names()
			roomNamesMsg := RoomsNamesMessage(names)
			cac.client.Send(roomNamesMsg)

			ujc := NewUserJoinedRoomMessage(cac.room, cac.client.ID())
			cac.client.Send(ujc)

		case roomName := <-ch.removeRoomRequests:

			if roomName == MainRoomName() {
				logger.Info("Cannot remove 'main' room")
				continue
			}

			delete(ch.rooms, roomName)
			msg := NewRemoveRoomMessage(roomName)
			ch.sendToEveryone(MainRoomName(), msg)

		case cac := <-ch.removeClientFromRoomRequest:
			logger.Infof("Remove client '%v' from room '%v'", cac.client, cac.room)

			room := ch.rooms[cac.room]
			room.RemoveClient(cac.client.ID())
			ujc := NewUserLeftRoomMessage(cac.room, cac.client.ID())
			cac.client.Send(ujc)

		case cac := <-ch.createRoomRequest:
			logger.Infof("Create room request from %v. Room name: %v", cac.client, cac.room)

			if !ch.roomNameValid(cac.room) {
				cac.client.Send(ErrorMessage("Invalid room name"))
				continue
			}

			if _, exists := ch.rooms[cac.room]; exists {
				logger.Infof("Room %v already exists. Client %v cannot create it", cac.room, cac.client)

				continue
			}

			// create new room with given name
			newRoom := NewRoom(cac.room, ch)
			newRoom.Start()
			newRoom.AddClient(cac.client)
			// add room to rooms' collection
			ch.rooms[cac.room] = newRoom

			ncm := NewCreateRoomMessage(cac.room)
			ch.sendToEveryone(MainRoomName(), ncm)

			ujc := NewUserJoinedRoomMessage(cac.room, cac.client.ID())
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

func (ch *Rooms) roomNameValid(name string) bool {
	if name == "" {
		logger.Info("invalid room name, name cannot be empty")
		return false
	}

	if !validRoomName.MatchString(name) {
		logger.Infof("invalid room name, name must match %v", roomNameRegexp)
	}

	return true
}

func (ch *Rooms) sendToEveryone(roomName string, msg *Message) {
	if room, ok := ch.rooms[roomName]; ok {
		logger.Infof("Send to room: %v", room.Name())
		room.SendToEveryone(msg)
	} else {
		logger.Infof("Cannot send message because the room %v doesn't exist", msg.Room)
	}
}

func (ch *Rooms) clientRooms(id string) []string {
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
func (ch *Rooms) CreateRoom(roomName string, client *Client) {
	ch.createRoomRequest <- clientAndRoom{
		client: client,
		room:   roomName,
	}
}

// RemoveClient removes client from all rooms.
func (ch *Rooms) RemoveClient(client *Client) {
	logger.Infof("Removing Client %v from all rooms", client)
	ch.removeClient <- client
}

// RemoveRoom removes room with given name.
func (ch *Rooms) RemoveRoom(roomName string) {
	ch.removeRoomRequests <- roomName
}

// ClientsRooms will return list of rooms to given client.
func (ch *Rooms) ClientsRooms(client *Client) {
	ch.roomsListRequests <- client
}

// AddClientToRoom adds given client to room with given name.
func (ch *Rooms) AddClientToRoom(roomName string, client *Client) {
	ch.addClientToRoomRequest <- clientAndRoom{
		client: client,
		room:   roomName,
	}
}

// RemoveClientFromRoom removes given client from room with given name.
func (ch *Rooms) RemoveClientFromRoom(roomName string, client *Client) {
	ch.removeClientFromRoomRequest <- clientAndRoom{
		client: client,
		room:   roomName,
	}
}

// SendMessageOnRoom sends given message to all clients of given room.
func (ch *Rooms) SendMessageOnRoom(message *Message) {
	ch.messageRequest <- message
}
