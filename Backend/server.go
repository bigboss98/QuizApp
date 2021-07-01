package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initializeDatabaseConnection(uri string) (*mongo.Client, error) {
	/*
	 * Initialize Database Connection to URI provided
	 *
	 * Param:
	 * 	-uri (string): URI where there is the database
	 */
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, _ := mongo.Connect(context.TODO(), clientOptions)

	// Check the connection
	error := client.Ping(context.TODO(), nil)

	if error != nil {
		fmt.Errorf("Invalid Database initialization")
	}
	return client, error
}

func decodeStartGameRequest(body io.ReadCloser) (Quiz, error) {
	/*
	 * Decode JSON provided by Start Game request
	 *
	 * Param:
	 * 	-body(io.ReadCloser):
	 */
	//var users []User
	var quiz Quiz
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	error := decoder.Decode(&quiz)

	quiz = quiz.setDefaultValues()
	return quiz, error
}

func decodeQuestionRequest(body io.ReadCloser) (Question, error) {
	/*
	 * Decode JSON provided by Insert Question request
	 *
	 * Param:
	 * 	-body(io.ReadCloser):
	 */
	var question Question
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	error := decoder.Decode(&question)
	if question.ID == "" {
		id, error := primitive.NewObjectIDFromTimestamp(time.Now()).MarshalJSON()
		if error != nil {
			log.Fatal(error)
		}
		question.ID = string(id[1 : len(id)-1])
	}
	return question, error
}
func startGame(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()
	verifyConnection(db, "QuizzoneDB")

	defer request.Body.Close()
	quiz, _ := decodeStartGameRequest(request.Body)
	insertQuizToDatabase(db, quiz)

	json_response, _ := json.MarshalIndent(quiz, "", "\t")
	response.Write([]byte(json_response))

}

func get_question(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
}

func answer_question(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

}

func updateQuestion(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()
	verifyConnection(db, "QuizzoneDB")

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
	verifyConnection(db, "QuizzoneDB")

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
	verifyConnection(db, "QuizzoneDB")

	defer request.Body.Close()
	questions := printQuestionsFromDatabase(db)
	json_response, _ := json.MarshalIndent(questions, "", "\t")
	response.Write([]byte(json_response))

}

func deleteQuestion(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()
	verifyConnection(db, "QuizzoneDB")

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
	verifyConnection(db, "QuizzoneDB")

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
	response.Header().Set("Content-Type", "application/json")
	pathParams := mux.Vars(request)
	db := openDatabase("QuizzoneDB")
	defer db.Close()
	verifyConnection(db, "QuizzoneDB")

	defer request.Body.Close()
	deleteQuizFromDatabase(db, pathParams["game_id"])

	response.Write([]byte("Emacs"))
}

func deleteQuizGames(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	db := openDatabase("QuizzoneDB")
	defer db.Close()
	verifyConnection(db, "QuizzoneDB")

	defer request.Body.Close()
	deleteQuizzesFromDatabase(db)

	response.Write([]byte("Emacs"))
}


func main() {
	router := mux.NewRouter()

	//Endpoints considered
	router.HandleFunc("/start_game", startGame).Methods("POST")                      //WORK
	router.HandleFunc("/insert_question", insertQuestion).Methods("POST")            //WORK
	router.HandleFunc("/update_question", updateQuestion).Methods("PUT")             //WORK make some test and choose what should be the response body 
	router.HandleFunc("/delete_question", deleteQuestion).Methods("DELETE")          //WORK change a little the response body
	router.HandleFunc("/delete_questions", deleteQuestions).Methods("DELETE")        //WORK change a little the response body
	router.HandleFunc("/get_question/{game_id}", get_question).Methods("GET")        //TO IMPLEMENT
	router.HandleFunc("/answer_question/{game_id}", answer_question).Methods("POST") // TO IMPLEMENT
	router.HandleFunc("/print_questions", printQuestions).Methods("GET")             //WORK change maybe a little response body 
	router.HandleFunc("/get_game/{game_id}", getGame).Methods("GET")                 //NOT IMPLEMENTED
	router.HandleFunc("/delete_game/{game_id}", deleteQuizGame).Methods("DELETE")    //WORK change response body
	router.HandleFunc("/delete_games", deleteQuizGames).Methods("DELETE")            //WORK change response body 
	log.Fatal(http.ListenAndServe(":8080", router))
}
