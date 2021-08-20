package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} //DEFINE Upgrader methods that upgrade HTTP request to Websocket

func decodeQuizGameRequest(body io.ReadCloser) (Quiz, error) {
	/*
	 * Decode Quiz JSON fields provided in Start Quiz request
	 *
	 * Param:
	 * 	-body(io.ReadCloser): Body object which contains part or all Quiz JSON fields
	 *
	 * Returns:
	 *   -quiz(Quiz): Quiz object derived from JSON fields, with missing fields set with initial values
	 *   -err(error): Error object with eventual errors from JSON decoding
	 */
	var quiz Quiz
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	error := decoder.Decode(&quiz)

	if error != nil {
		log.Fatal(error)
	}
	//quiz = quiz.setInitialValues() //Initialize with initial values other quiz fields
	return quiz, error
}

func decodeQuestionRequest(body io.ReadCloser) (Question, error) {
	/*
	 * Decode JSON provided by Insert Question request
	 *
	 * Param:
	 * 	-body(io.ReadCloser): body which contains part or all JSON fields of a Question
	 *
	 * Returns:
	 *   -question(Question): question object which contains all decoded Question fields
	 *   -err(error): error object used to identify error happens during JSON fields decoding
	 */
	var question Question
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	error := decoder.Decode(&question)

	// Case when you insert a new Question and you don't have an question ID
	if question.ID == "" {
		id, error := primitive.NewObjectIDFromTimestamp(time.Now()).MarshalJSON()
		if error != nil {
			log.Printf("%s", error)
		}
		question.ID = string(id[1 : len(id)-1])
	}
	return question, error
}

func decodeAnswerGiven(body io.ReadCloser) (AnsweredQuestion, error) {
	/*
	 * Decode Answered Question JSON fields provided in Answer Question request
	 *
	 * Param:
	 * 	-body(io.ReadCloser): Body object which contains part or all Quiz JSON fields
	 *
	 * Returns:
	 *   -answer(AnsweredQuestion): AnsweredQuestion object derived from JSON fields
	 *   -err(error): Error object with eventual errors from JSON decoding
	 */
	var answer AnsweredQuestion
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	error := decoder.Decode(&answer)

	if error != nil {
		log.Printf("%s", error)
	}

	answer_id, error := primitive.NewObjectIDFromTimestamp(time.Now()).MarshalJSON()
	if error != nil {
		log.Printf("%s", error)
	}
	answer.ID = string(answer_id)
	return answer, error
}

func signUp(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	var user User
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&user)

	if err != nil {
		log.Fatal(err)
	}

	user.Password, err = GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	db := openDatabase("QuizzoneDB")
	defer db.Close()

	insertUserToDatabase(db, user)
	json_user := encodeUser(user, "\t", "")
	response.Write([]byte(json_user))
}

func signIn(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	var authdetails Authentication
	err := json.NewDecoder(request.Body).Decode(&authdetails)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	authuser := getUserFromDatabase(db, authdetails.Name)
	check := CheckPasswordHash(authdetails.Password, authuser.Password)

	if !check {
		log.Printf("Error: Username or Password is incorrect")
	}

	validToken, err := GenerateJWT(authuser.Name, authuser.Role)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	var token Token
	token.Name = authuser.Name
	token.Role = authuser.Role
	token.TokenString = validToken

	json_token := encodeSignIn(token, "\t", "")
	response.Write([]byte(json_token))

}

func startConn(wsServer *WsServer, response http.ResponseWriter, request *http.Request) {
	/*
	 * API to create and "start" a Conn given Users names in request body
	 *
	 * Params:
	 * -wsServer(*WsServer): Websocket connection
	 * -response(http.ResponseWriter): response object used to give created Quiz as response
	 * -request(*http.Request): request object with all information needed to process and create
	 *							a new Quiz object on DB
	 */
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	requestVars := mux.Vars(request)
	conn, err := upgrader.Upgrade(response, request, nil)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Printf("Error in Websocket upgrades with error %s", err)
	}

	db := openDatabase("QuizzoneDB")
	defer db.Close()
	userDb := getUserFromDatabase(db, requestVars["name"])
	user := newUser(conn, wsServer, userDb.Name, userDb.Password, userDb.Role, "Not Authorized")
	go user.writeMessage()
	go user.readMessage()

	wsServer.register <- user

}

func updateQuestion(response http.ResponseWriter, request *http.Request) {
	/*
	 * API to update a Question on DB
	 * -response(http.ResponseWriter): response object used to give updated question as response
	 * -request(*http.Request): request object with all information needed to process and update
	 *							Question object on DB
	 */
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	updateQuestionToDatabase(db, question)
	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func insertQuestion(response http.ResponseWriter, request *http.Request) {
	/*
	 * API used to insert a new question on DB
	 * Params:
	 * -response(http.ResponseWriter): response object used to give created Question as response
	 * -request(*http.Request): request object with all information needed to process and create
	 *							a new Question object on DB
	 */
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	insertQuestionToDatabase(db, question)
	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func deleteQuestion(response http.ResponseWriter, request *http.Request) {
	/*
	 * API used to delete Question from DB
	 * Params:
	 * -response(http.ResponseWriter): response object used to give deleted Question as response
	 * -request(*http.Request): request object with all information needed to process and delete
	 							a Question object on DB
	*/
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	deleteQuestionFromDatabase(db, question.ID)

	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func deleteQuestions(response http.ResponseWriter, request *http.Request) {
	/*
	 * API used to delete all Question from DB
	 * Params:
	 * -response(http.ResponseWriter): response object used to give all deleted Question as response
	 * -request(*http.Request): request object with all information needed to process and delete
	 							all Question objects on DB
	*/
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	deleteQuestionsFromDatabase(db)
	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func deleteQuizGame(response http.ResponseWriter, request *http.Request) {
	/*
	 * API used by endpoint delete_quizGame, which delete a quiz game object
	 * from DB given its ID.
	 *
	 * Param:
	 * -response(http.ResponseWriter): response object used to write response for this API
	 * -request(*http.Request): request object which contains JSON body and parameters needed
	 *                          to process and deleted Quiz games
	 */
	response.Header().Set("Content-Type", "application/json")
	pathParams := mux.Vars(request)
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	deleteQuizFromDatabase(db, pathParams["game_id"])

	response.Write([]byte("Emacs"))
}

func deleteQuizGames(response http.ResponseWriter, request *http.Request) {
	/*
	 * API used by endpoint delete_quizGames, which delete a all quiz games object from DB.
	 *
	 * Param:
	 * -response(http.ResponseWriter): response object used to write response for this API
	 * -request(*http.Request): request object which contains JSON body and parameters needed
	 *                          to process and deleted Quiz games
	 */
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	deleteQuizzesFromDatabase(db)

	response.Write([]byte("Emacs"))
}

func main() {
	/*
	 * Main function used to define and manage API endpoints
	 */
	router := mux.NewRouter()

	wsServer := NewWebsocketServer()
	go wsServer.Run()
	//Endpoints considered
	router.HandleFunc("/start_quiz/{name}", func(w http.ResponseWriter, r *http.Request) {
		startConn(wsServer, w, r)
	})
	router.HandleFunc("/sign_up", signUp).Methods("POST") //WORK only create User table on the DB
	router.HandleFunc("/sign_in", signIn).Methods("POST") //Create User table on the DB 
	router.HandleFunc("/insert_question", insertQuestion).Methods("POST", "OPTIONS")  //WORK
	router.HandleFunc("/update_question", updateQuestion).Methods("PUT")              //WORK make some test and choose what should be the response body
	router.HandleFunc("/delete_question", deleteQuestion).Methods("DELETE")           //WORK change a little the response body
	router.HandleFunc("/delete_questions", deleteQuestions).Methods("DELETE")         //WORK change a little the response body       //WORK
	router.HandleFunc("/delete_quizGame/{game_id}", deleteQuizGame).Methods("DELETE") //WORK change response body
	router.HandleFunc("/delete_allQuizGames", deleteQuizGames).Methods("DELETE")      //WORK change response body
	log.Fatal(http.ListenAndServe(":8080", router))
}
