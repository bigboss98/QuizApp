/*
 * Model file contains all struct and functions to model Data model for
 * Quiz APP (User, Question and Quiz game object)
 */

package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	/*
	 * Set Default values for a Quiz object
	 *
	 *
	 */
	quiz.Scores = make([]int, len(quiz.Users))
	quiz.NumPlayers = len(quiz.Users)
	quiz.Questions = choiceQuestions(num_question)
	quiz.Status = "started"
	game_id, error := primitive.NewObjectIDFromTimestamp(time.Now()).MarshalJSON()
	if error != nil {
		log.Fatal(error)
	}
	quiz.Game_ID = string(game_id[1 : len(game_id)-1])
	return quiz
}

func choiceQuestions(num_question int) []Question {
	/*
	 * Choices randomly n Questions to use as question of a quiz game
	 *
	 * param:
	 * -num_question(int): number of question to choose
	 */
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
