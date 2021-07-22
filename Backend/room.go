package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Room struct {
	/*
	 * Represents a Room where Game starts
	 * Fields:
	 * -Name(string): Name of the Room
	 * -users(map[*User]bool): list of users associated to the room
	 * -register(chan *User): channel where there will be client to be registered to the room
	 * -unregister(chan *User): channel where there will be client to be unregistered from the room
	 * -broadcast(chan *Message): message to be broadcasted to all users
	 * -ready(map[*User]bool): map of clients with their ready game status
	 * -status(string): indicates whether game is started, ended or not already started
	 */
	Name       string `json:name`
	users      map[*User]bool
	register   chan *User
	unregister chan *User
	broadcast  chan *Message
	ready      map[*User]bool
	quiz       *Quiz
	status     string
}

func NewRoom(name string) *Room {
	/*
	 * Creates a new room object and return it
	 * Params:
	 * -name(string): Name of Room object
	 */
	return &Room{
		Name:       name,
		users:      make(map[*User]bool),
		register:   make(chan *User),
		unregister: make(chan *User),
		broadcast:  make(chan *Message),
		ready:      make(map[*User]bool),
		quiz:       nil,
		status:     "not started",
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {
	for {
		select {

		case user := <-room.register:
			room.registerUserInRoom(user)

		case user := <-room.unregister:
			room.unregisterUserInRoom(user)

		case message := <-room.broadcast:
			switch message.Action {
				case StartGameAction:
					room.notifyStartMessage(message)

			//case GetQuestionAction:
			//room.notifyGetQuestion(message)

				case AnswerQuestionAction:
					room.notifyAnswerQuestion(message)
				//room.broadcastToClientsInRoom(message.encode())
			}
		}

	}
}

func (room *Room) registerUserInRoom(user *User) {
	/*
	 * Register user in a room with its status set to "not ready"
	 * Params:
	 * -user(*User): User object to be registered on the room
	 */
	room.notifyUserJoined(user)
	room.users[user] = true
	room.ready[user] = false
}

func (room *Room) unregisterUserInRoom(user *User) {
	/*
		     * Unregister user from a Room object
			 * -user(*User): user to be unregistered from the room
	*/
	if _, ok := room.users[user]; ok {
		delete(room.users, user)
		delete(room.ready, user)
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	/*
	 * Broadcast Message to all clients of the room
	 * Params:
	 * -message([]byte): message to be sended to all users of the room
	 */
	for user := range room.users {
		user.send <- message
	}
}

func (room *Room) notifyAnswerQuestion(message *Message) {
	var answer AnsweredQuestion
	err := json.Unmarshal(message.Message, &answer)
	if err != nil {
		log.Printf("Error in Answer Question encoding: %s", err)
	}
	room.quiz.answerQuestion(&answer)
	json_response := encodeAnswerQuestion(answer, "\t", "")
	response := &Response{
		Action:  message.Action,
		Message: string(json_response),
		Sender:  message.Sender,
		Target:  message.Target,
	}
	room.broadcastToClientsInRoom(response.encode())
}

func (room *Room) notifyUserJoined(user *User) {
	/*
	 * Notify Join Message to all users of the Room with action set to JoinRoomAction
	 * Params:
	 * -user(*User): User object where notify that joins the room
	 */
	response := &Response{
		Action:  JoinRoomAction,
		Target:  room,
		Message: "{" + fmt.Sprintf(welcomeMessage, user.GetName()) + "}",
		Sender:  user,
	}

	room.broadcastToClientsInRoom(response.encode())
}

func (room *Room) notifyStartMessage(message *Message) {
	/*
	 * Notify Start message to all users of the Room when action is set to StartGameAction
	 * Params:
	 * -room(*Room): room object responsable to notify start message
	 * -message(*Message): message to be notified to all users of the room
	 */
	var response_message string
	if room.status == "started" {
		room.quiz = createQuiz(room.users)
		question := getCurrentQuestion(room.quiz)
		response_message = encodeQuestion(&question)

	} else {
		response_message = "{" + "statusRoom: " + room.status + ", userStatus: ready}"
	}
	response := &Response{
		Action:  message.Action,
		Message: response_message,
		Sender:  message.Sender,
		Target:  message.Target,
	}
	room.broadcastToClientsInRoom(response.encode())
}

func (room *Room) GetName() string {
	/*
	 * Return Name of the room object
	 */
	return room.Name
}
