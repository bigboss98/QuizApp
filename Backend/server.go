package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	quiz = quiz.setInitialValues() //Initialize with initial values other quiz fields
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
			log.Fatal(error)
		}
		question.ID = string(id[1 : len(id)-1])
	}
	return question, error
}

func startQuiz(response http.ResponseWriter, request *http.Request) {
	/*
	 * API to create and "start" a Quiz given Users names in request body
	 *
	 * Params:
	 * -response(http.ResponseWriter): response object used to give created Quiz as response
	 * -request(*http.Request): request object with all information needed to process and create
	 							a new Quiz object on DB
	*/
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	quiz, _ := decodeQuizGameRequest(request.Body)
	insertQuizToDatabase(db, quiz)

	json_response, _ := json.MarshalIndent(quiz, "", "\t")
	response.Write([]byte(json_response))

}

func get_question(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	pathParams := mux.Vars(request)
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	quiz := getQuizFromDatabase(db, pathParams["game_id"])
	question := quiz.getCurrentQuestion()

	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func answer_question(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	pathParams := mux.Vars(request)
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	quiz := getQuizFromDatabase(db, pathParams["game_id"])
	question := quiz.getCurrentQuestion()
	
	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))

}

func updateQuestion(response http.ResponseWriter, request *http.Request) {
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
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	insertQuestionToDatabase(db, question)
	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func printQuestions(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	questions := printQuestionsFromDatabase(db)
	json_response, _ := json.MarshalIndent(questions, "", "\t")
	response.Write([]byte(json_response))

}

func deleteQuestion(response http.ResponseWriter, request *http.Request) {
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
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	deleteQuestionsFromDatabase(db)
	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func getGame(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

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
	router := mux.NewRouter()

	//Endpoints considered
	router.HandleFunc("/start_quiz", startQuiz).Methods("POST")                       //WORK
	router.HandleFunc("/insert_question", insertQuestion).Methods("POST")             //WORK
	router.HandleFunc("/update_question", updateQuestion).Methods("PUT")              //WORK make some test and choose what should be the response body
	router.HandleFunc("/delete_question", deleteQuestion).Methods("DELETE")           //WORK change a little the response body
	router.HandleFunc("/delete_questions", deleteQuestions).Methods("DELETE")         //WORK change a little the response body
	router.HandleFunc("/get_question/{game_id}", get_question).Methods("GET")         //WORK change a little the response body when question are ended 
	router.HandleFunc("/answer_question/{game_id}", answer_question).Methods("POST")  // TO IMPLEMENT
	router.HandleFunc("/print_questions", printQuestions).Methods("GET")              //WORK change maybe a little response body
	router.HandleFunc("/get_quizGame/{game_id}", getGame).Methods("GET")              //NOT IMPLEMENTED -> only declare endpoint manages on server.go
	router.HandleFunc("/delete_quizGame/{game_id}", deleteQuizGame).Methods("DELETE") //WORK change response body
	router.HandleFunc("/delete_allQuizGames", deleteQuizGames).Methods("DELETE")      //WORK change response body
	log.Fatal(http.ListenAndServe(":8080", router))
}
