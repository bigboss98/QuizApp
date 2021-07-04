/*
 * Model file contains all struct and functions to model Data model for
 * Quiz APP (User, Question and Quiz game object)
 */

package main

import (
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const num_question int = 10

type User struct {
	Name       string `json:name`
	QuizPlayed []Quiz `json:quiz_played`
}

type Question struct {
	ID       string   `json:id`
	Question string   `json:question`
	Language string   `json:language`
	Author   string   `json:author`
	Choices  []string `json:choices`
	Category string   `json:category`
	Answer   string   `json:answer`
}

type AnsweredQuestion struct {
	ID             string   `json:id`
	User           []User   `json:user`
	Answers        []string `json:answers`
	Correct_Answer string   `json:correct_answer`
}

type Quiz struct {
	Game_ID          string     `json:game_ID`
	Users            []string   `json:users`
	Winner           string     `json:winner`
	Scores           []int      `json:scores`
	Status           string     `json:status`
	Questions        []string `json:questions`
	AnswerGiven 	 []string   `json:answer_given`
	NumPlayers  	 int        `json:num_players`
	current_question int //Index for current question on Questions field
}

func (quiz Quiz) setInitialValues() Quiz {
	/*
	 * Set Default values for a Quiz object
	 *
	 *
	 */
	quiz.Scores = make([]int, len(quiz.Users))
	quiz.NumPlayers = len(quiz.Users)
	quiz.Questions = transformQuestionsToString(choiceQuestions(num_question))
	quiz.Status = "started"
	quiz.current_question = 0
	game_id, error := primitive.NewObjectIDFromTimestamp(time.Now()).MarshalJSON()
	if error != nil {
		log.Fatal(error)
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
	rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })

	if len(questions) > num_question {
		return questions[:num_question]
	}
	return questions
}

func transformQuestionsToString(questions []Question) []string{
	var questions_id []string 
	for _, question := range questions{
		questions_id = append(questions_id, question.ID)
	}
	return questions_id
}

func (quiz Quiz) getCurrentQuestion() Question{
	var question Question
	current_position := quiz.current_question
	if current_position < len(quiz.Questions){
		question_id := quiz.Questions[quiz.current_question]
		db := openDatabase("QuizzoneDB")
		defer db.Close()
		
		quiz.current_question = quiz.current_question + 1 
		question = getQuestionFromDatabase(db, question_id)
		updateQuizToDatabase(db, quiz)
	}
	return question
}

/*
func transformStringToQuestion(questions_id []string) []Question{

}
*/