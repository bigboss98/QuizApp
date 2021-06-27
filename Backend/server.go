package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const num_question int = 40

type User struct {
	Name string `json:name`
}

type Question struct {
	ID       string   `json:id`
	Question string   `json:question`
	Choices  []string `json:choices`
	Category string   `json:category`
	Answer   string   `json:answer`
}

type Quiz struct {
	Game_ID    string     `json:game_ID`
	Users      []User     `json:users`
	Winner     string     `json:winner`
	Scores     []int      `json:scores`
	Status     string     `json:status`
	Questions  []Question `json:questions`
	NumPlayers int        `json:num_players`
}

func (quiz Quiz) setDefaultValues() Quiz {
	quiz.Scores = make([]int, len(quiz.Users))
	quiz.NumPlayers = len(quiz.Users)
	quiz.Questions = make([]Question, num_question)
	quiz.Status = "started"
	game_id, error := primitive.NewObjectIDFromTimestamp(time.Now()).MarshalJSON()
	if error != nil {
		log.Fatal(error)
	}
	quiz.Game_ID = string(game_id[1 : len(game_id)-1])
	return quiz
}

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
	//quiz = { , quiz.Users, "", "not started", }
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
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	//pathParams := mux.Vars(request)

	// get collection as ref
	collection := client.Database("quizdb").Collection("quiz")

	defer request.Body.Close()
	game, _ := decodeStartGameRequest(request.Body)
	_, _ = collection.InsertOne(context.TODO(), game)

	json_response, _ := json.MarshalIndent(game, "", "\t")
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
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	//pathParams := mux.Vars(request)

	// get collection as ref
	collection := client.Database("quizdb").Collection("question")

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	filter := bson.M{"id": bson.M{"$eq": "60d7a9f7cfc3ef2d8fe2735d"}}

	update := bson.D{
		{"$set", bson.D{{"question", question.Question},
			{"choices", question.Choices},
			{"category", question.Category},
			{"answer", question.Answer}}},
	}
	_, _ = collection.UpdateOne(context.TODO(), filter, update)

	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func insertQuestion(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	//pathParams := mux.Vars(request)

	// get collection as ref
	collection := client.Database("quizdb").Collection("question")

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	_, _ = collection.InsertOne(context.TODO(), question)

	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func printQuestions(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	//pathParams := mux.Vars(request)

	// get collection as ref
	collection := client.Database("quizdb").Collection("question")

	defer request.Body.Close()
	cursor, _ := collection.Find(context.TODO(), bson.D{})

	var question Question
	var questions []Question
	for cursor.Next(context.TODO()) {
		cursor.Decode(&question)
		questions = append(questions, question)
	}
	json_response, _ := json.MarshalIndent(questions, "", "\t")
	response.Write([]byte(json_response))
}

func deleteQuestion(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	//pathParams := mux.Vars(request)

	// get collection as ref
	collection := client.Database("quizdb").Collection("question")

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	//filter := bson.M{"id": bson.M{"$eq": "60d7a9f7cfc3ef2d8fe2735d"}}

	_, _ = collection.DeleteOne(context.TODO(), bson.D{{}})

	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}
func main() {
	router := mux.NewRouter()

	//Endpoints considered
	router.HandleFunc("/start_game", startGame).Methods("POST")           //WORK
	router.HandleFunc("/insert_question", insertQuestion).Methods("POST") //WORK
	router.HandleFunc("/update_question/", updateQuestion).Methods("PUT") // WORK change a little response body 
	router.HandleFunc("/delete_question/", deleteQuestion).Methods("DELETE") //WORK change a little the response body
	router.HandleFunc("/get_question/{game_id}", get_question).Methods("GET")        //TO IMPLEMENT
	router.HandleFunc("/answer_question/{game_id}", answer_question).Methods("POST") // TO IMPLEMENT
	router.HandleFunc("/print_questions", printQuestions).Methods("GET")             //WORK
	log.Fatal(http.ListenAndServe(":8080", router))
}
