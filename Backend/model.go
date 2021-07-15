/*
 * Model file contains all struct and functions to model Data model for
 * Quiz APP (User, Question, Answered Question and Quiz game object)
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const num_question int = 10

type User struct {
	/*
	 * Represent a User of the game
	 * Fields:
	 * -Name(string): name of the player 
	 */
	Name       string `json:name`
	//createdAt time.Time 
}

type Question struct {
	/*
	 * Represent a Question of Quiz game
	 * Fields:
	 * -ID(string): an ID string to uniquely identify a question
	 * -Question(string): question text
	 * -Language(string): language of the question
	 * -Author(string): Name of Author user 
	 * -Choices([]string): array of string that represent answer's choices
	 * -Category(string): category of the question
	 * -Answer(string): correct Answer of the question
	 */
	ID       string   `json:id`
	Question string   `json:question`
	Language string   `json:language`
	Author   string   `json:author`
	Choices  []string `json:choices`
	Category string   `json:category`
	Answer   string   `json:"answer, omitempty"`
	//createdAt time.Time  
}

type AnsweredQuestion struct {
	/*
	 * Represent an Answered Question of Quiz
	 * Fields:
	 * -ID(string): identification string of the AnsweredQuestion
	 * -Question(string): ID of question to identify question
	 * -Users([]string): users Names of the answered question
	 * -Answers([]string): answers choices 
	 * -CorrectAnswer(string): correct answer of the question
	 * -TimeToAnswer(time.Time): time required to Answer 
	 */
	ID             string    `json:id`
	Question       string    `json:question`
	Users          []string  `json:users`
	Answers        []string  `json:answers`
	Correct_Answer string    `json:correct_answer`
	Score		   int       `json:score`
	TimeToAnswer   int       `json:TimeToAnswer`
}

type Quiz struct {
	/*
	 * Represent a Quiz object
	 * Fields:
	 * -Game_ID(string): ID of Quiz
	 * -Users([]string): Names of Quiz users
	 * -Winner(string): Name of Quiz winner
	 * -Scores(string): Scores of all quiz players
	 * -Status(string): Status of the game
	 * -Questions([]string): Array of Question IDs
	 * -AnswerGiven([]string): Array of AnsweredQuestion IDs
	 * -NumPlayers(int): number of Quiz players
	 * -current_question(int): current question index to retrieve current question
	 * -createdAt(time.Time): Time where Quiz match was created 
	 */
	Game_ID          string   `json:game_ID`
	Users            []string `json:users`
	Winner           string   `json:winner`
	Scores           []int    `json:scores`
	Status           string   `json:status`
	Questions        []string `json:questions`
	AnswerGiven      []string `json:answer_given`
	NumPlayers       int      `json:num_players`
	current_question int      //Index for current question on Questions field
	createdAt        time.Time
}

const maxTimeAnswer = 20; 
const basePoint = 200;
const maxIncrement = 100; 
const decrementPoint = -200; 

func (quiz Quiz) setInitialValues() Quiz {
	/*
	 * Set Default values for a Quiz object which has the following default values:
	 * current_question = 0
	 * Status = "started"
	 * NumPlayers = len(quiz.Users)
	 * Question = id of n randomly chosen Questions
	 * createdAt = time.Now() 
	 */
	quiz.Scores = make([]int, len(quiz.Users))
	quiz.NumPlayers = len(quiz.Users)
	quiz.Questions = transformQuestionsToString(choiceQuestions(num_question))
	quiz.Status = "started"
	quiz.current_question = 0
	quiz.createdAt = time.Now()
	game_id, error := primitive.NewObjectIDFromTimestamp(time.Now()).MarshalJSON()
	if error != nil {
		log.Printf("Error: %s", error)
	}
	quiz.Game_ID = string(game_id[1 : len(game_id)-1]) //Need to remove first char since it is a \ char
	return quiz
}

func choiceQuestions(num_question int) []Question {
	/*
	 * Choices randomly n Questions to use as question of a quiz game
	 *
	 * param:
	 * -num_question(int): number of question to choose
	 */
	db := openDatabase("QuizzoneDB")
	defer db.Close()

	questions := printQuestionsFromDatabase(db)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(questions), 
	             func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })

	if len(questions) > num_question {
		return questions[:num_question]
	}
	return questions
}

func transformQuestionsToString(questions []Question) []string {
	/*
	 * Transform array of n questions in un array of n Question ID
	 * used to save in Database and to represent questions field in Quiz object
	 */
	var questions_id []string
	for _, question := range questions {
		questions_id = append(questions_id, question.ID)
	}
	return questions_id
}

func getCurrentQuestion(db *pgxpool.Pool, quiz *Quiz) Question {
	/*
	 * Return current QUestion of a Quiz object, updating current question value in Quiz object
	 * Params:
	 * -db(*pgxpool.Pool):database connection used to retrieve question given question ID
	 * -quiz(*Quiz): Quiz object used to retrieve question ID of current question 
	 *
	 * Return a Question object that represent the current question or default empty question
	 * when all questions of a Quiz are already retrieved
	 */
	var question Question
	current_position := quiz.current_question
	if current_position < len(quiz.Questions) {
		question_id := quiz.Questions[quiz.current_question]
		quiz.current_question = quiz.current_question + 1
		question = getQuestionFromDatabase(db, question_id)
		updateQuizToDatabase(db, *quiz)
	}
	return question
}

func (quiz Quiz) answerQuestion(db *pgxpool.Pool, question Question, answer *AnsweredQuestion) {
	/*
	 * Answer Question used to check answers given and to insert an AnsweredQuestion on DB
	 * Params:
	 * -quiz(Quiz): quiz object where we want to answer a question
	 * -question(Question): question object that we wnat to answer
	 * -answer(*AnsweredQuestion): AnsweredQuestion that represent answer to question
	 */
	checkAnswerGiven(answer.Answers, question.Choices)

	quiz.AnswerGiven = append(quiz.AnswerGiven, answer.ID)
	answer.Question = question.ID
	answer.Correct_Answer = question.Answer
	answer.Users = quiz.Users
	answer.Score = answer.computeScores()
	quiz.Scores[0] = quiz.Scores[0] + answer.Score 
	updateQuizToDatabase(db, quiz)
	insertAnswerToDatabase(db, *answer)
}

func (answer AnsweredQuestion) computeScores() int {
	if answer.Answers[0] == answer.Correct_Answer {
		fmt.Sprint(answer.TimeToAnswer)
		return basePoint + maxIncrement / maxTimeAnswer * 
						   (maxTimeAnswer - answer.TimeToAnswer)
	}
	return decrementPoint
}

func checkAnswerGiven(answers []string, choices []string) error {
	/*
	 * Check whether given answers are valid choices for Question answer
	 * Param:
	 * -answers([]string): array of answer given by users
	 * -choices([]string): array of Question answer choices
	 *
	 * Return an error when an answer is not a valid answer choice for the question
	 */
	var index int
	for _, answer := range answers {
		found := false
		for index = 0; index < len(choices) && !found; index++ {
			if answer == choices[index] {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("answer %s is not a valid choice for the question", answer)
		}
	}
	return nil
}

func encodeInitialGame(quiz Quiz, indent string, prefix string) []byte {
	/*
	 * Encode StartQuiz response body 
	 * Params:
	 * -quiz(Quiz): Quiz object 
	 * -indent(string): indent pattern used to indent JSON response
	 * -prefix(string): prefix pattern used to indent JSON response
	 *
	 * Return a JSON object represented by a []byte with the following structure
	 * {
	 *		"game_id": ID of Quiz game
	 * }
	 */
	encode_quiz := map[string]interface{}{
		"game_id": quiz.Game_ID,
	}
	json_question, _ := json.MarshalIndent(encode_quiz, prefix, indent)
	return json_question
}

func encodeAnswerQuestion(answer AnsweredQuestion, indent string, prefix string) []byte {
	/*
	 * Encode Answer Question response body 
	 * Params:
	 * -answer(AnsweredQuestion): answeredQuestion object 
	 * -indent(string): indent pattern used to indent JSON response
	 * -prefix(string): prefix pattern used to indent JSON response
	 *
	 * Return a JSON object represented by a []byte with the following structure
	 * {
	 *		"correct_answer": correct answer of the AnsweredQuestion
	 *      "guess": return whether user guess the correct answer of the Question
	 * }
	 */
	guess_answer := answer.Answers[0] == answer.Correct_Answer
	encode_answer := map[string]interface{}{
		"correct_answer": answer.Correct_Answer,
		"guess": guess_answer,
		"score": answer.Score,
	}
	json_answer, _ := json.MarshalIndent(encode_answer, prefix, indent)
	return json_answer
}

