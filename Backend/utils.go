package main 

/*
 * Contains all const, common variable shared in all source code 
 */
const num_question int = 3
const maxTimeAnswer = 20
const basePoint = 200
const maxIncrement = 100
const decrementPoint = -200
const welcomeMessage = "%s joined the room"

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const SendMessageAction = "send-message"
const AnswerQuestionAction = "answer-question"
const GetQuestionAction = "get-question"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const StartGameAction = "start-game"
const UserLeftAction = "user-left"
const UserJoinedAction = "user-join"
const EndGameAction = "end-game"
const GetResultAction = "get-results"
const GetPlayersAction = "get-players"
const AuthentificationAction = "auth"