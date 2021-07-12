package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/lib/pq"
)

const (
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
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
}

func openDatabase(dbName string) *pgxpool.Pool {
	/*
	 * Open a database connection
	 * Param:
	 * -dbName(string): name of database
	 * Return a sql.DB object which represent a pool of DB connections
	 */
	db, err := pgxpool.Connect(context.Background(), dsn(dbName))
	if err != nil {
		log.Printf("Error %s when opening DB %s", err, dbName)
		return nil
	}
	log.Printf("Connection to %s happened correctly", dbName)
	//db.SetMaxOpenConns(20)                 //Number of max open connections available on database dbName
	//db.SetMaxIdleConns(20)                 //Max number of idle connections available on database dbName
	//db.SetConnMaxLifetime(time.Minute * 5) //Maximum time for a database connection

	return db
}

func insertQuestionToDatabase(db *pgxpool.Pool, question Question) {
	/*
	 *
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
	 *
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
	 *
	 */
	deleteQuestion := `delete from "question" where id=$1`
	_, err := db.Exec(context.Background(), deleteQuestion, question_id)

	if err != nil {
		log.Printf("Errors %s deleting Question to DB", err)
	}
	log.Printf("Delete question with ID %s successfully", question_id)
}

func deleteQuestionsFromDatabase(db *pgxpool.Pool) {
	/*
	 *
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
	 *
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
	return questions
}

func getQuestionFromDatabase(db *pgxpool.Pool, question_id string) Question {
	/*
	 * Get a Question from Database given question ID
	 *
	 * Params:
	 * -db(*sql.DB): database connection
	 * -question_id(string): ID of Question to be retrieved from DB
	 *
	 * Return Retrieved question or Log error
	 */
	selectQuestion := `SELECT * FROM "question" WHERE id=$1`
	log.Printf("ID Question: %s", question_id)
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

	return Question{id, question, language, author, choices, category, answer}

}

//func insertUserToDatabase(db *sql.DB, user User){

//}

func insertAnswerToDatabase(db *pgxpool.Pool, answer AnsweredQuestion){
	insertAnswer := `insert into "answeredquestion" values($1, $2, $3, $4)`
	_, err := db.Exec(context.Background(), insertAnswer, answer.ID, answer.Answers, answer.Question,
					  "{" + strings.Join(answer.Users, ",") + "}")

	if err != nil {
		log.Printf("Errors %s inserting Answer Question to DB", err)
	}
	log.Printf("Insert Answer Question to DB successfully")
}

func getQuizFromDatabase(db *pgxpool.Pool, quiz_id string) Quiz {
	/* 
	 *
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

	err := row.Scan(&id, &users, &winner, &status, &num_players, &questions,
		pq.Array(&answer_given), &scores, &current_question)

	if err != nil {
		log.Printf("Errors %s in retrieved Quiz from DB", err)
	}

	if winner == nil {
		temp := ""
		winner = &temp
	}
	return Quiz{id, users, *winner, scores, status, questions, answer_given, num_players, current_question}

}

func insertQuizToDatabase(db *pgxpool.Pool, quiz Quiz) {
	/*
	 *
	 */
	insertQuiz := `insert into "quiz"(quiz_id, users, status, scores, num_players, questions, current_question) values($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(context.Background(), insertQuiz, quiz.Game_ID, "{"+strings.Join(quiz.Users, ",")+"}",
		quiz.Status, pq.Array(quiz.Scores), quiz.NumPlayers, "{"+strings.Join(quiz.Questions, ",")+"}",
		quiz.current_question)

	if err != nil {
		log.Printf("Errors %s insering Question to DB", err)
	}
	log.Printf("Insert quiz to DB successfully")
}

func updateQuizToDatabase(db *pgxpool.Pool, quiz Quiz) {
	/*
	 *
	 */
	var err error
	if quiz.Winner != "" {
		updateQuestion := `update "quiz" SET users=$2, winner=$3, status=$4, scores=$5, questions=$6, num_players=$7,
						   current_question=$8, answers_given=$9 where quiz_id=$1`
		_, err = db.Exec(context.Background(), updateQuestion, quiz.Game_ID, quiz.Users, quiz.Winner, quiz.Status,
			quiz.Scores, "{"+strings.Join(quiz.Questions, ",")+"}", quiz.NumPlayers, quiz.current_question,
			pq.Array(quiz.AnswerGiven))

	} else {
		updateQuestion := `update "quiz" SET users=$2, status=$3, scores=$4, questions=$5, num_players=$6,
						   current_question=$7, answers_given=$8 where quiz_id=$1`
		_, err = db.Exec(context.Background(), updateQuestion, quiz.Game_ID, quiz.Users, quiz.Status,
			quiz.Scores, "{"+strings.Join(quiz.Questions, ",")+"}", quiz.NumPlayers, quiz.current_question,
			pq.Array(quiz.AnswerGiven))

	}

	if err != nil {
		log.Printf("Errors %s when updating Quiz to DB", err)
	}
	log.Printf("Updated quiz game to DB successfully")
}

func deleteQuizFromDatabase(db *pgxpool.Pool, quiz_id string) {
	/*
	 *
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
	 *
	 */
	deleteQuestion := `delete from "quiz"`
	_, err := db.Exec(context.Background(), deleteQuestion)

	if err != nil {
		log.Printf("Errors %s deleting all Quiz game from DB", err)
	}
	log.Printf("Deleted all quiz game from DB successfully")
}
