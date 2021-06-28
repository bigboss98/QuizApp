package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
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
	quiz.Questions = generateQuestions(num_question)
	quiz.Status = "started"
	game_id, error := primitive.NewObjectIDFromTimestamp(time.Now()).MarshalJSON()
	if error != nil {
		log.Fatal(error)
	}
	quiz.Game_ID = string(game_id[1 : len(game_id)-1])
	return quiz
}

func generateQuestions(num_question int) []Question {
	client, _ := initializeDatabaseConnection("mongodb://localhost")

	// get collection as ref
	collection := client.Database("quizdb").Collection("question")

	cursor, _ := collection.Find(context.TODO(), bson.D{})

	var question Question
	var questions []Question
	for cursor.Next(context.TODO()) {
		cursor.Decode(&question)
		questions = append(questions, question)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })

	if len(questions) > num_question {
		return questions[:num_question]
	}
	return questions
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
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	pathParams := mux.Vars(request)

	collection := client.Database("quizdb").Collection("quiz")
	defer request.Body.Close()
	filter := bson.M{"game_id": bson.M{"$eq": pathParams["game_id"]}}

	var quiz Quiz
	_ = collection.FindOne(context.TODO(), filter).Decode(&quiz)
	question := quiz.Questions[0]
	quiz.Questions = quiz.Questions[1:]
	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func answer_question(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	//client, _ := initializeDatabaseConnection("mongodb://localhost")
	//pathParams := mux.Vars(request)

	defer request.Body.Close()
	//question, _ := decodeQuestionRequest(request.Body)
	//filter := bson.M{"id": bson.M{"$eq": pathParams["game_id"]}}

}

func updateQuestion(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	//pathParams := mux.Vars(request)

	// get collection as ref
	collection := client.Database("quizdb").Collection("question")

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	filter := bson.M{"id": bson.M{"$eq": question.ID}}

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
	filter := bson.M{"id": bson.M{"$eq": question.ID}}

	_, _ = collection.DeleteOne(context.TODO(), filter)

	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func deleteQuestions(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	//pathParams := mux.Vars(request)

	// get collection as ref
	collection := client.Database("quizdb").Collection("question")

	defer request.Body.Close()
	question, _ := decodeQuestionRequest(request.Body)
	//filter := bson.M{"id": bson.M{"$eq": "60d7a9f7cfc3ef2d8fe2735d"}}

	_, _ = collection.DeleteMany(context.TODO(), bson.D{{}})

	json_response, _ := json.MarshalIndent(question, "", "\t")
	response.Write([]byte(json_response))
}

func getGame(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	pathParams := mux.Vars(request)

	collection := client.Database("quizdb").Collection("quiz")

	filter := bson.M{"id": bson.M{"$eq": pathParams["game_id"]}}
	var quiz Quiz
	_ = collection.FindOne(context.TODO(), filter).Decode(&quiz)

	json_response, _ := json.MarshalIndent(quiz, "", "\t")
	response.Write([]byte(json_response))
}

func deleteQuizGame(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	pathParams := mux.Vars(request)
	collection := client.Database("quizdb").Collection("quiz")

	filter := bson.M{"id": bson.M{"$eq": pathParams["game_id"]}}

	_, _ = collection.DeleteOne(context.TODO(), filter)

	response.Write([]byte("Delete a Quiz Game"))
}

func deleteQuizGames(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	client, _ := initializeDatabaseConnection("mongodb://localhost")
	collection := client.Database("quizdb").Collection("quiz")

	_, _ = collection.DeleteMany(context.TODO(), bson.D{{}})

	response.Write([]byte("Emacs"))
}

func main() {
	router := mux.NewRouter()

	//Endpoints considered
	router.HandleFunc("/start_game", startGame).Methods("POST")                      //WORK set also Question
	router.HandleFunc("/insert_question", insertQuestion).Methods("POST")            //WORK
	router.HandleFunc("/update_question", updateQuestion).Methods("PUT")             // WORK change a little response body
	router.HandleFunc("/delete_question", deleteQuestion).Methods("DELETE")          //WORK change a little the response body
	router.HandleFunc("/delete_questions", deleteQuestions).Methods("DELETE")        //WORK change a little the response body
	router.HandleFunc("/get_question/{game_id}", get_question).Methods("GET")        //WORK make some test to be sure about -->Insert an update to the question
	router.HandleFunc("/answer_question/{game_id}", answer_question).Methods("POST") // TO IMPLEMENT
	router.HandleFunc("/print_questions", printQuestions).Methods("GET")             //WORK but also insert query choices
	router.HandleFunc("/get_game/{game_id}", getGame).Methods("GET")                 //WORK
	router.HandleFunc("/delete_game/{game_id}", deleteQuizGame).Methods("DELETE")    //WORK
	router.HandleFunc("/delete_games", deleteQuizGames).Methods("DELETE")            //WORK maybe not return nothing
	log.Fatal(http.ListenAndServe(":8080", router))
}
