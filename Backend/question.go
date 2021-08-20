package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

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
	ID             string   `json:id`
	Question       string   `json:question`
	User           string   `json:user`
	Answer         string   `json:answer`
	Correct_Answer string   `json:correct_answer`
	Score          int      `json:score`
	TimeToAnswer   int      `json:TimeToAnswer`
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

func getCurrentQuestion(quiz *Quiz) Question {
	/*
	 * Return current QUestion of a Quiz object, updating current question value in Quiz object
	 * Params:
	 * -quiz(*Quiz): Quiz object used to retrieve question ID of current question
	 *
	 * Return a Question object that represent the current question or default empty question
	 * when all questions of a Quiz are already retrieved
	 */
	var question Question
	current_position := quiz.current_question
	if current_position < len(quiz.Questions) {
		question = quiz.Questions[current_position]
		quiz.current_question = quiz.current_question + 1
		question.Answer = ""
	}
	return question
}

func (answer AnsweredQuestion) computeScores() int {
	/*
	 * Compute Score of an Answer with the following rules:
	 * 1) if Answer is wrong the score is given by decrement point constant(by default -200)
	 * 2) If answer is correct the score is basePoint(default +200) + 
	 			(maxIncrement /maxTimeAnswer) * (maxTimeAnswer-userTimeToAnswer)
	 */
	if answer.Answer == answer.Correct_Answer {
		return basePoint + maxIncrement/maxTimeAnswer*
			(maxTimeAnswer-answer.TimeToAnswer)
	}
	return decrementPoint
}

func checkAnswerGiven(answer string, choices []string) error {
	/*
	 * Check whether given answers are valid choices for Question answer
	 * Param:
	 * -answer(string): answer given by user
	 * -choices([]string): array of Question answer choices
	 *
	 * Return an error when an answer is not a valid answer choice for the question
	 */
	var index int
	for index = 0; index < len(choices); index++ {
		if answer == choices[index] {
			return nil 
			}
	}
	return fmt.Errorf("answer %s is not a valid choice for the question", answer)
}

func encodeQuestion(question *Question) string {
	/*
	 * Encode Question response message with the following format:
	 * {
	 *		"question": JSON object that represent the question
	 *      "status": indicates whether the question is retrieved or not from a Quiz match
	 * }
	 */
	var status bool = true
	if question.Question == "" {
		status = false //Status used to show whether we retrieve or not question from a Quiz match
	}
	encode_question := map[string]interface{}{
		"question": question,
		"status":   status,
	}
	json_answer, _ := json.MarshalIndent(encode_question, "", "\t")
	return string(json_answer)
}


