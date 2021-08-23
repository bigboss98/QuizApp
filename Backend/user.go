package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)



type User struct {
	/*
	 * Represent a User of the game
	 * Fields:
	 * -Name(string): name of the player
	 * -Status(string): status of the player("ready", "not ready", "Not Authorized")
	 * -Password(string): Password of the player (Saved in DB with only their Hash)
	 * -conn(*websocket.Conn): connection object of an User
	 * -room(*Room): represent Room object of the user 
	 * -wsServer(*WsServer): WebSocket server object where all connections are been directed
	 * -send(chan []byte): JSON message sended by User
	 */
	Name     string `json:name`
	Status   string `json:status`
	Password string `json:password`
	Role     string `json:role`
	conn     *websocket.Conn
	room     *Room
	wsServer *WsServer
	send     chan []byte

	//createdAt time.Time
}

func newUser(conn *websocket.Conn, wsServer *WsServer, name string, password string, role string, status string) *User {
	/*
	 * Creates a New User with also its Connection
	 *
	 * Params:
	 *	-conn(*websocket.Conn): websocket connection 
	 *  -wsServer(*WsServer): Websocket object associated to user
	 *  -name(string): Name of User 
	 *  -password(string): Password of the User 
	 */
	return &User{
		conn:     conn,
		Name:     name,
		Status: status,
		Password: password,
		wsServer: wsServer,
		room:     nil,
		send:     make(chan []byte, 3000),
	}
}

func (user *User) writeMessage() {
	/*
	 * Write message in a concurrent way and notify other user on Websocket connection
	 * Params:
	 * -user(*User): user object that has writeMessage operation 
	 */
	//ticker := time.NewTicker(pingPeriod)
	defer func() {
		//ticker.Stop()
		user.conn.Close()
	}()
	for {
		select {
		case message, ok := <- user.send:
			//client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				user.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := user.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(message)

			// Attach queued chat messages to the current websocket message.
			length := len(user.send)
			for index := 0; index < length; index++ {
				writer.Write(newline)
				writer.Write(<-user.send)
			}

			if err := writer.Close(); err != nil {
				return
			}
			//case <-ticker.C:
			//	user.conn.SetWriteDeadline(time.Now().Add(writeWait))
			//	if err := user.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			//		return
			//	}
		}
	}
}

func (user *User) readMessage() {
	/*
	 * Read Message in a concurrent way received by a User from Websocket connection
	 * When receive a message call handleNewMessage to handle the possible message action
	 *
	 */
	defer func() {
		user.disconnect()
	}()

	//user.conn.SetReadLimit(maxMessageSize)
	//user.conn.SetReadDeadline(time.Now().Add(pongWait))
	//client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := user.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
		} else {
			user.handleNewMessage(jsonMessage)
		}
	}
}

func (user *User) handleNewMessage(jsonMessage []byte) {
	/*
	 * Handle New message received from an User 
	 * Params:
	 * -user(*User): user object that receive a new message
	 * -jsonMessage([]byte): JSON message received from Websocket
	 */
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	// Attach the client object as the sender of the messsage.
	message.Sender = user

	switch message.Action {
		case EndGameAction:
			user.handleGame(&message)

		case AnswerQuestionAction:
			user.handleGame(&message)

		case GetQuestionAction:
			user.handleGame(&message)

		case StartGameAction:
			user.handleStartGameMessage(&message)

		case JoinRoomAction:
			user.handleJoinRoomMessage(message)

		case LeaveRoomAction:
			user.handleLeaveRoomMessage(message)

		case GetPlayersAction:
			user.handleGetPlayersMessage(&message)

		case AuthentificationAction:
			user.handleAuthentification(&message)
	}
}

func (user *User) handleGetPlayersMessage(message *Message) {
	if user.Status != "Not Authorized" {
		roomName := message.Target.Name 

		room := user.wsServer.findRoomByName(roomName)
	
		room.broadcast <- message
	} 

}

func (user *User) handleAuthentification(message *Message) {
	log.Printf("Request Authentification")
	if IsAuthorized(message.Token) {
		user.Status = "not ready"
		user.send <- []byte("User Authorized")
	}else {
		log.Printf("User Not Authorized")
		close(user.send)
		user.conn.Close()
	}
}

func (user *User) handleGame(message *Message) {
	if user.Status != "Not Authorized" {
		roomName := message.Target.Name

		room := user.wsServer.findRoomByName(roomName)

		room.broadcast <- message
	}else{
		log.Printf("NOT AUTHORIZED")
	}
}

func (user *User) handleStartGameMessage(message *Message) {
	/* 
	 * Handle StartGame Message which set status of user in a room to ready
	 * When all user status are set to ready it really start the game and provides 
	 * the first question and write on database 
	 *
	 * Params:
	 * -user(*User): user object that want start the game 
	 * -message(Message): message Object 
	 */
	
	if user.Status != "Not Authorized" { 
		roomName := message.Target.Name

		room := user.wsServer.findRoomByName(roomName)

		if room != nil && room.status != "started"{ 
			room.users[user.GetName()].Status = "ready"
			room.ready[user] = true 

			var startedGame = "started" 
			for _, status := range room.ready {
				if startedGame == "started" && !status {
					startedGame = "not started" 
				}
			}
			room.status = startedGame
			room.broadcast <- message
		}
	} 
}

func (user *User) handleJoinRoomMessage(message Message) {
	/*
	 * Handle Join a room of the user
	 * Params:
	 * -user(*User): user Object which joins the room 
	 * -message(Message): message received from User to join Room 
	 */
	if user.Status != "Not Authorized" {
		roomName := message.Target.Name
		log.Printf("Name Room: %s", roomName)

		room := user.wsServer.findRoomByName(roomName)
		if room == nil {
			room = user.wsServer.createRoom(roomName)
		}
		user.room = room
		room.register <- user
	}else{
		log.Printf(user.Status)
		log.Printf("NOT AUTHORIZED")
	}
}

func (user *User) handleLeaveRoomMessage(message Message) {
	/*
	 * Handle Leave from a room of the user
	 * Params:
	 * -user(*User): user Object which leaves the room 
	 * -message(Message): message received from User to leave Room 
	 */
	if user.Status != "Not Authorized" {
		room := user.wsServer.findRoomByName(message.Target.Name)
		if room != nil && user.room != nil && user.room == room {
			user.room = nil
			room.unregister <- user
		}
	}
}

func (user *User) disconnect() {
	/*
	 * Disconnect an User from Room and websocket
	 * Params:
	 * -user(*User): User object where we want to disconnect from its websocket connection 
	 */
	user.wsServer.unregister <- user
	user.room.unregister <- user
	close(user.send)
	user.conn.Close()
}

func (user *User) GetName() string {
	/*
	 * Return Name of User
	 * Params:
	 * -user: User object which we return its name 
	 */
	return user.Name
}

func getPlayers(players map[string]*User) []User {
	var users []User
	for _, user := range players {
		if user != nil {
			users = append(users, *user)
		}
	}
	return users 
}

func encodePlayers(users []User, room *Room) string {
	encode_players := map[string]interface{}{
		"players": users,
	}
	json_answer, _ := json.MarshalIndent(encode_players, "", "\t")
	return string(json_answer)
}

func encodeUser(user User, indent string, prefix string) string {
	encode_user := map[string]interface{}{
		"name": user.Name, 
		"role": user.Role,
	}
	json_user, _ := json.MarshalIndent(encode_user, prefix, indent)
	return string(json_user)
}