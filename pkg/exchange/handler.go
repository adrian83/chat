package exchange

import (
	"fmt"
)

type Handler interface {
	Handle(msg *Message) error
}

func NewRoute(msgType string, handler Handler) *Route {
	return &Route{
		msgType: msgType,
		handler: handler,
	}
}

type Route struct {
	msgType string
	handler Handler
}

func (r *Route) Handle(msg *Message) error {
	if r == nil {
		return fmt.Errorf("cannot handle message: %v", msg)
	}

	return r.handler.Handle(msg)
}

func (r *Route) MsgType() string {
	return r.msgType
}

func NewRouter() *Router {
	return &Router{
		routes: map[string]*Route{},
	}
}

type Router struct {
	routes map[string]*Route
}

func (r *Router) RegisterRoute(route *Route) {
	r.routes[route.MsgType()] = route
}

func (r *Router) FindRoute(msgType string) *Route {
	return r.routes[msgType]
}

// ----

func NewAddClientToRoomHandler(rooms *Rooms, client *Client) *AddClientToRoomHandler {
	return &AddClientToRoomHandler{
		rooms:  rooms,
		client: client,
	}
}

type AddClientToRoomHandler struct {
	rooms  *Rooms
	client *Client
}

func (h *AddClientToRoomHandler) Handle(msg *Message) error {
	h.rooms.AddClientToRoom(msg.Room, h.client)
	return nil
}

// ----

func NewSendMsgToRoomHandler(rooms *Rooms) *SendMsgToRoomHandler {
	return &SendMsgToRoomHandler{
		rooms: rooms,
	}
}

type SendMsgToRoomHandler struct {
	rooms *Rooms
}

func (h *SendMsgToRoomHandler) Handle(msg *Message) error {
	h.rooms.SendMessageOnRoom(msg)
	return nil
}

// ----

func NewCreateRoomHandler(rooms *Rooms, client *Client) *CreateRoomHandler {
	return &CreateRoomHandler{
		rooms:  rooms,
		client: client,
	}
}

type CreateRoomHandler struct {
	rooms  *Rooms
	client *Client
}

func (h *CreateRoomHandler) Handle(msg *Message) error {
	h.rooms.CreateRoom(msg.Room, h.client)
	return nil
}

// ----

func NewRemoveClientFromRoomHandler(rooms *Rooms, client *Client) *RemoveClientFromRoomHandler {
	return &RemoveClientFromRoomHandler{
		rooms:  rooms,
		client: client,
	}
}

type RemoveClientFromRoomHandler struct {
	rooms  *Rooms
	client *Client
}

func (h *RemoveClientFromRoomHandler) Handle(msg *Message) error {
	h.rooms.RemoveClientFromRoom(msg.Room, h.client)
	return nil
}

// ----

func NewLogoutHandler(client *Client) *LogoutHandler {
	return &LogoutHandler{
		client: client,
	}
}

type LogoutHandler struct {
	client *Client
}

func (h *LogoutHandler) Handle(msg *Message) error {
	h.client.Send(msg)
	return nil
}
