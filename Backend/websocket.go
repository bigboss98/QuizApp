package main

import (
	"encoding/json"
	"log"
)

type WsServer struct {
	/*
	 * Represent all WebSocket connections with
	 * all available rooms
	 *
	 * Params:
	 * -users(map[*User]bool): Map of Clients associated to WsSocket instance 
	 * -register(chan *User): channel where there will be user to be registed in WsServer
	 * -unregister(chan *User): channel where there will be user to be removed from WsServer instance
	 * -broadcast(chan []byte): channel where there is message to be broadcasted to all users 
	 * -rooms(map[*Room]bool): Map of Rooms associated to WsSocket instance
	 */
	users    map[*User]bool
	register   chan *User
	unregister chan *User
	broadcast  chan []byte
	rooms      map[*Room]bool
}

type Message struct {
	/*
	 * Indicates the input Message given by Websocket's user 
	 * Fields:
	 * -Action(string): indicates the action of the message (possible action are provided in common.go)
	 * -Message(AnsweredQuestion): Answered Question needed only with AnsweredQuestion otherwise is not considered 
	 * -Target(*Room): indicates the room target of the message 
	 * -Sender(*User): indicates the sender of the message 
	 */
	Action  string           `json:"action"`
	Token   string           `json:"token"`
	Message AnsweredQuestion `json:"message"`
	Target  *Room            `json:"target"`
	Sender  *User            `json:"sender"`
}

type Response struct {
	/*
	 * Response message given after every websocket communication
	 * Fields:
	 * -Action(string): return the type of message useful as information
	 * -Message(string): information given as response to a 
	 * -Target(*Room): Room target of Response message -> Makes sense as response??
	 * -Sender(*User): Sender of the message --> Makes Sense?
	 */
	Action  string           `json:action`
	Message string           `json:message`
	Target  *Room            `json:target`
	Sender  *User            `json:sender`
}

func (message *Message) encode() []byte {
	/*
	 * Encode a Message JSON object from a Message object 
	 */
	jsonMessage, err := json.Marshal(message)

	if err != nil {
		log.Println(err)
	}

	return jsonMessage
}

func (response *Response) encode() []byte {
	/*
	 * Encode a Response object as JSON from Response object 
	 */
	jsonMessage, err := json.Marshal(response)

	if err != nil {
		log.Println(err)
	}

	return jsonMessage
}

func NewWebsocketServer() *WsServer{
	return &WsServer{
		users: make(map[*User]bool),
		register: make(chan *User),
		unregister: make(chan *User),
		broadcast: make(chan []byte),
		rooms: make(map[*Room]bool),				
	}
}

func (server *WsServer) findRoomByName(name string) *Room {
	for room := range server.rooms {
		if room.GetName() == name {
			return room
		}
	}

	return nil 
}

func (server *WsServer) createRoom(name string) *Room {
	room := NewRoom(name)
	go room.RunRoom()
	server.rooms[room] = true

	return room
}

func (server *WsServer) Run() {
	log.Println("Started RUN")
	for {
		select {

			case user := <-server.register:
				log.Printf("%s is being registered", user.Name)
				server.registerClient(user)

			case user := <-server.unregister:
				server.unregisterClient(user)

			case message := <-server.broadcast:
				server.broadcastToClients(message)
		}

	}
}

func (server *WsServer) registerClient(user *User) {
	server.notifyClientJoined(user)
	//server.listOnlineClients(client)
	server.users[user] = true
}

func (server *WsServer) unregisterClient(user *User) {
	if _, ok := server.users[user]; ok {
		delete(server.users, user)
		server.notifyClientLeft(user)
	}
}

func (server *WsServer) notifyClientJoined(user *User) {
	response := &Response{
		Action: UserJoinedAction,
		Sender: user,
	}

	server.broadcastToClients(response.encode())
}

func (server *WsServer) notifyClientLeft(user *User) {
	response := &Response{
		Action: UserLeftAction,
		Sender: user,
	}

	server.broadcastToClients(response.encode())
}

func (server *WsServer) broadcastToClients(message []byte) {
	for user := range server.users {
		user.send <- message
	}
}
