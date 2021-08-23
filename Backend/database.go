package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
)

const (
	/*
	 * Connection constant used to connect to database
	 */
	user     = "bigboss98"
	password = "Xyz?^(?yZ02"
	host     = "localhost"
	port     = 5432
)

func dsn(dbName string) string {
	/*
	 * Return string name used to provide host, port, user, password needed to access database
	 * Param:
	 * -dbName(string): Name of database
	 */
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
}

func openDatabase(dbName string) *pgxpool.Pool {
	/*
	 * Open a database connection returning a pgxpool.Pool object that represent
	 * a Pool of Database connections
	 * Param:
	 * -dbName(string): name of database
	 */
	db, err := pgxpool.Connect(context.Background(), dsn(dbName))

	if err != nil {
		log.Printf("Error %s when opening DB %s", err, dbName)
		return nil
	}
	log.Printf("Connection to %s happened correctly", dbName)

	return db
}

func insertQuestionToDatabase(db *pgxpool.Pool, question Question) {
	/*
	 * Insert Question Object to DB
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 * -question(Question): question object to be inserted
	 */
	insertQuestion := `insert into "question" values($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(context.Background(), insertQuestion, question.ID, question.Question, question.Language,
		question.Category, question.Answer, question.Author, "{"+strings.Join(question.Choices, ",")+"}")

	if err != nil {
		log.Printf("Errors %s insering Question to DB", err)
	}
	log.Printf("Insert question to DB successfully")
}

func updateQuestionToDatabase(db *pgxpool.Pool, question Question) {
	/*
	 * Update Question object to DB
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 * -question(Question): question object to be updated on DB
	 */
	updateQuestion := `update "question" SET question=$2, language=$3, category=$4, answer=$5, author=$6, choices=$7 where id=$1`
	_, err := db.Exec(context.Background(), updateQuestion, question.ID, question.Question, question.Language, question.Category,
		question.Answer, question.Author, "{"+strings.Join(question.Choices, ",")+"}")

	if err != nil {
		log.Printf("Errors %s when updating Question to DB", err)
	}
	log.Printf("Updated question to DB successfully")
}
func deleteQuestionFromDatabase(db *pgxpool.Pool, question_id string) {
	/*
	 * Delete question object with a given Question ID from DB
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 * -question_id(string): Question ID of question to be deleted from DB
	 */
	deleteQuestion := `delete from "question" where id=$1`
	_, err := db.Exec(context.Background(), deleteQuestion, question_id)

	if err != nil {
		log.Printf("Errors %s deleting Question to DB", err)
	}
	log.Printf("Deleted question with ID %s", question_id)
}

func deleteQuestionsFromDatabase(db *pgxpool.Pool) {
	/*
	 * Delete all Questions from Database
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 */
	deleteQuestion := `delete from "question"`
	_, err := db.Exec(context.Background(), deleteQuestion)

	if err != nil {
		log.Printf("Errors %s deleting all Questions from DB", err)
	}
	log.Printf("Deleted all question from DB successfully")
}

func printQuestionsFromDatabase(db *pgxpool.Pool) []Question {
	/*
	 * Retrieve all Questions from DB returning an array of Question object
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 */
	rows, err := db.Query(context.Background(), `SELECT * FROM "question"`)

	if err != nil {
		log.Printf("Errors %s when printing all Questions from DB", err)
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var id string
		var question string
		var language string
		var category string
		var answer string
		var author string
		var choices []string
		err = rows.Scan(&id, &question, &language, &category, &answer,
			&author, pq.Array(&choices))

		if err != nil {
			log.Printf("Errors %s in Scanning retrieved Questions from DB", err)
		}
		questions = append(questions, Question{id, question, language, author, choices, category, answer})
	}
	log.Printf("Retrieved %d questions from Database", len(questions))
	return questions
}

func getQuestionFromDatabase(db *pgxpool.Pool, question_id string) Question {
	/*
	 * Get a Question from Database given question ID
	 *
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connection
	 * -question_id(string): ID of Question to be retrieved from DB
	 *
	 * Return Retrieved question object
	 */
	selectQuestion := `SELECT * FROM "question" WHERE id=$1`
	row := db.QueryRow(context.Background(), selectQuestion, question_id)
	var id string
	var question string
	var language string
	var category string
	var answer string
	var author string
	var choices []string

	err := row.Scan(&id, &question, &language, &category, &answer, &author, pq.Array(&choices))

	if err != nil {
		log.Printf("Errors %s in retrieved Questions from DB", err)
	}
	log.Printf("Retrieved Question from DB with ID: %s", question_id)

	return Question{id, question, language, author, choices, category, answer}
}

//func insertUserToDatabase(db *sql.DB, user User){

//}

func insertAnswerToDatabase(db *pgxpool.Pool, answer AnsweredQuestion) {
	/*
	 * Insert Answer object to DB
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 * -answer(AnsweredQuestion): answeredQuestion object to be inserted on DB
	 */
	//insertAnswer := `insert into "answeredquestion" values($1, $2, $3, $4, $5, $6)`
	/*
	_, err := db.Exec(context.Background(), insertAnswer, answer.ID, answer.Answers, answer.Question,
		"{"+strings.Join(answer.Users, ",")+"}", answer.TimeToAnswer, answer.Score)

	if err != nil {
		log.Printf("Errors %s inserting Answer Question to DB", err)
	}
	log.Printf("Insert Answer Question to DB successfully")
	*/
}

func getQuizFromDatabase(db *pgxpool.Pool, quiz_id string) Quiz {
	/*
	 * Retrieve quiz object from DB given its Quiz ID
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 * -quiz_id(string): ID of Quiz object to be retrieved
	 */
	getQuiz := `select * from "quiz" where quiz_id=$1`
	row := db.QueryRow(context.Background(), getQuiz, quiz_id)

	var id string
	var users []string
	var winner *string
	var status string
	var scores []int
	var num_players int
	var questions []string
	var answer_given []string
	var current_question int
	var created_at time.Time
	err := row.Scan(&id, &users, &winner, &status, &num_players, &questions,
		pq.Array(&answer_given), &scores, &current_question, &created_at)

	if err != nil {
		log.Printf("Errors %s in retrieved Quiz from DB", err)
	}
	log.Printf("Retrieved Quiz object with ID %s", quiz_id)

	if winner == nil { //CASE with no winner and by default we set with ""
		temp := ""
		winner = &temp
	}
//	return Quiz{id, users, *winner, scores, status, questions, answer_given, num_players, current_question, created_at}
	return Quiz{}
}

func insertQuizToDatabase(db *pgxpool.Pool, quiz Quiz) {
	/*
	 * Insert Quiz object on DB
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 * -quiz(Quiz): quiz object to be inserted
	 */
	//insertQuiz := `insert into "quiz"(quiz_id, users, status, scores, num_players, questions,
	//			   current_question, created_at) values($1, $2, $3, $4, $5, $6, $7, $8)`
	/*
	_, err := db.Exec(context.Background(), insertQuiz, quiz.Game_ID, "{"+strings.Join(quiz.Users, ",")+"}",
		quiz.Status, pq.Array(quiz.Scores), quiz.NumPlayers, "{"+strings.Join(quiz.Questions, ",")+"}",
		quiz.current_question, quiz.createdAt)

	if err != nil {
		log.Printf("Errors %s inserting Question to DB", err)
	}
	log.Printf("Insert quiz to DB successfully")
	*/
}

func updateQuizToDatabase(db *pgxpool.Pool, quiz Quiz) {
	/*
	 * Update Quiz object on DB given its Quiz id
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 * -quiz(Quiz): quiz object to be updated on DB
	 */
	//var err error
	/*
	//Case with no winner already
	if quiz.Winner != "" {
		updateQuestion := `update "quiz" SET users=$2, winner=$3, status=$4, scores=$5, questions=$6, num_players=$7,
						   current_question=$8, answers_given=$9, created_at=$10 where quiz_id=$1`
		_, err = db.Exec(context.Background(), updateQuestion, quiz.Game_ID, quiz.Users, quiz.Winner, quiz.Status,
			quiz.Scores, "{"+strings.Join(quiz.Questions, ",")+"}", quiz.NumPlayers, quiz.current_question,
			pq.Array(quiz.AnswerGiven), quiz.createdAt)

	} else { //Case of END of game with winner
		updateQuestion := `update "quiz" SET users=$2, status=$3, scores=$4, questions=$5, num_players=$6,
						   current_question=$7, answers_given=$8, created_at=$9 where quiz_id=$1`
		_, err = db.Exec(context.Background(), updateQuestion, quiz.Game_ID, quiz.Users, quiz.Status,
			quiz.Scores, "{"+strings.Join(quiz.Questions, ",")+"}", quiz.NumPlayers, quiz.current_question,
			pq.Array(quiz.AnswerGiven), quiz.createdAt)

	}

	if err != nil {
		log.Printf("Errors %s when updating Quiz to DB", err)
	}
	log.Printf("Updated quiz game to DB successfully")
	*/
}

func deleteQuizFromDatabase(db *pgxpool.Pool, quiz_id string) {
	/*
	 * Delete Quiz object from database given its ID
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 * -quiz_id(string): ID of Quiz object to be deleted from DB
	 */
	deleteQuestion := `delete from "quiz" where quiz_id=$1`
	_, err := db.Exec(context.Background(), deleteQuestion, quiz_id)

	if err != nil {
		log.Printf("Errors %s deleting Quiz to DB", err)
	}
	log.Printf("Deleted quiz with ID %s successfully", quiz_id)
}

func deleteQuizzesFromDatabase(db *pgxpool.Pool) {
	/*
	 * Delete all Quiz matches from database
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connections
	 */
	deleteQuestion := `delete from "quiz"`
	_, err := db.Exec(context.Background(), deleteQuestion)

	if err != nil {
		log.Printf("Errors %s deleting all Quiz game from DB", err)
	}
	log.Printf("Deleted all quiz game from DB successfully")
}

func insertUserToDatabase(db *pgxpool.Pool, user User) {
	/*
	 * Insert Used to DB during Sign Up 
	 * 
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connection used to insert User object 
	 * -user(User): User object to be inserted on DB
	 */
	insertUser := `insert into "user"(name, password, role) values ($1, $2, $3)`
	_, err := db.Exec(context.Background(), insertUser, user.Name, user.Password, user.Role)

	if err != nil {
		log.Printf("Error: %s", err)
	}
	log.Printf("Insert User in Database")

}

func getUserFromDatabase(db *pgxpool.Pool, name string) *User{
	/*
	 * Get User from DB used for Login 
	 * 
	 * Params:
	 * -db(*pgxpool.Pool): Pool of DB connection used to retrieve DB 
	 * -name(string): name of User
	 */
	getUser := `select * from "user" where name=$1`
	row := db.QueryRow(context.Background(), getUser, name)
	var username string 
	var password string 
	var role string 
	err := row.Scan(&username, &role, &password)

	if err != nil {
		log.Printf("Error: %s", err)
	}
	log.Printf("Get User from Database")
	
	return &User{
		Name: username,
		Password: password,
		Role: role,
		Status: "not ready",
	}
}