/*
 * Model file contains all struct and functions to model Data model for
 * Quiz APP (User, Question, Answered Question and Quiz game object)
 */

package main

import (
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	Game_ID          string         `json:game_ID`
	Users            map[*User]bool `json:users`
	Winner           string         `json:winner`
	Scores           map[string]int  `json:scores`
	Status           string         `json:status`
	Questions        []Question     `json:questions`
	AnswerGiven      []string       `json:answer_given`
	NumPlayers       int            `json:num_players`
	current_question int            //Index for current question on Questions field
	createdAt        time.Time
}

func createQuiz(users map[*User]bool) *Quiz {
	/*
	 * Creates a Quiz object from a map of Users 
	 */
	game_id, err := primitive.NewObjectIDFromTimestamp(time.Now()).MarshalJSON()
	if err != nil {
		log.Printf("Error: %s", err)
	}
	return &Quiz{
		Scores:           make(map[string]int),
		NumPlayers:       len(users),
		Status:           "started",
		Users:            users,
		Questions:        choiceQuestions(num_question),
		current_question: 0,
		createdAt:        time.Now(),
		Game_ID:          string(game_id[1 : len(game_id)-1]), //Need to remove first char since there is a \ char
	}
}

func (quiz Quiz) answerQuestion(answer *AnsweredQuestion) {
	/*
	 * Answer Question used to check answers given and to insert an AnsweredQuestion on DB
	 * Params:
	 * -quiz(Quiz): quiz object where we want to answer a question
	 * -question(Question): question object that we wnat to answer
	 * -answer(*AnsweredQuestion): AnsweredQuestion that represent answer to question
	 */
	question := quiz.Questions[quiz.current_question-1]
	checkAnswerGiven(answer.Answer, question.Choices)

	quiz.AnswerGiven = append(quiz.AnswerGiven, answer.ID)
	answer.Question = question.ID
	answer.Correct_Answer = question.Answer
	//answer.Users = quiz.Users
	log.Print(answer.Answer)
	log.Printf("Correct Answer: %s", answer.Correct_Answer)
	answer.Score = answer.computeScores()
	quiz.Scores[answer.User] = quiz.Scores[answer.User] + answer.Score
	//updateQuizToDatabase(db, quiz)
	//insertAnswerToDatabase(db, *answer)
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
	guess_answer := answer.Answer == answer.Correct_Answer
	encode_answer := map[string]interface{}{
		"correct_answer": answer.Correct_Answer,
		"guess":          guess_answer,
		"score":          answer.Score,
	}
	json_answer, _ := json.MarshalIndent(encode_answer, prefix, indent)
	return json_answer
}

func encodeGetQuiz(quiz Quiz, indent string, prefix string) []byte {
	/*
	 * Encode Get Quiz response body
	 * Params:
	 * -quiz(Quiz): Quiz object
	 * -indent(string): indent pattern used to indent JSON response
	 * -prefix(string): prefix pattern used to indent JSON response
	 *
	 * Return a JSON object represented by a []byte with the following structure
	 * {
	 *		"score": score of Player on Game
	 *      "status": status of Quiz game
	 *      "guess": return whether user guess the correct answer of the Question
	 * }
	 */
	encode_answer := map[string]interface{}{
		"score":  quiz.Scores,
		"status": quiz.Status,
		"winner": quiz.Winner,
	}
	json_answer, _ := json.MarshalIndent(encode_answer, prefix, indent)
	return json_answer
}

func (quiz *Quiz) endGame() *Quiz {
	/*
	 * End the game and return the new Quiz object 
	 * 
	 */
	if quiz.current_question == num_question && len(quiz.AnswerGiven) == num_question {
		quiz.Status = "ended"
		winner := ""
		bestScore := -500000
		for user, score := range quiz.Scores {
			if score > bestScore {
				bestScore = score
				winner = user 
			}
		}
		quiz.Winner = winner 
	}
	return quiz
}
