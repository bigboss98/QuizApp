package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

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

func openDatabase(dbName string) *sql.DB {
	/*
	 * Open a database connection
	 * Param:
	 * -dbName(string): name of database
	 * Return a sql.DB object which represent a pool of DB connections
	 */
	db, err := sql.Open("postgres", dsn(dbName))
	if err != nil {
		log.Printf("Error %s when opening DB %s", err, dbName)
		return nil
	}
	db.SetMaxOpenConns(20)                 //Number of max open connections available on database dbName
	db.SetMaxIdleConns(20)                 //Max number of idle connections available on database dbName
	db.SetConnMaxLifetime(time.Minute * 5) //Maximum time for a database connection

	return db
}

func verifyConnection(db *sql.DB, dbName string) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	err := db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbName)
}

func insertQuestionToDatabase(db *sql.DB, question Question) {
	insertQuestion := `insert into "question" values($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(insertQuestion, question.ID, question.Question, question.Language,
		question.Category, question.Answer, question.Author, "{"+strings.Join(question.Choices, ",")+"}")

	if err != nil {
		log.Printf("Errors %s insering Question to DB", err)
	}
	log.Printf("Insert question to DB successfully")
}

func updateQuestionToDatabase(db *sql.DB, question Question) {
	updateQuestion := `update "question" SET question=$2, language=$3, category=$4, answer=$5, author=$6, choices=$7 where id=$1`
	_, err := db.Exec(updateQuestion, question.ID, question.Question, question.Language, question.Category,
		question.Answer, question.Author, "{"+strings.Join(question.Choices, ",")+"}")

	if err != nil {
		log.Printf("Errors %s when updating Question to DB", err)
	}
	log.Printf("Updated question to DB successfully")
}
func deleteQuestionFromDatabase(db *sql.DB, question_id string) {
	deleteQuestion := `delete from "question" where id=$1`
	_, err := db.Exec(deleteQuestion, question_id)

	if err != nil {
		log.Printf("Errors %s deleting Question to DB", err)
	}
	log.Printf("Delete question with ID %s successfully", question_id)
}

func deleteQuestionsFromDatabase(db *sql.DB) {
	deleteQuestion := `delete from "question"`
	_, err := db.Exec(deleteQuestion)

	if err != nil {
		log.Printf("Errors %s deleting all Questions from DB", err)
	}
	log.Printf("Deleted all question from DB successfully")
}

func printQuestionsFromDatabase(db *sql.DB) []Question {
	rows, err := db.Query(`SELECT * FROM "question"`)

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

//func insertUserToDatabase(db *sql.DB, user User){

//}

/*
func getQuizFromDatabase(db *sql.DB, quiz_id string) Quiz{
	getQuiz := `select from "quiz" where id=$1`
	row := db.QueryRow(getQuiz, quiz_id)

	if err != nil {
		log.Printf("Errors %s deleting Question to DB", err)
	}
	log.Printf("Delete question with ID %s successfully", question_id)
}
*/
func insertQuizToDatabase(db *sql.DB, quiz Quiz) {
	insertQuiz := `insert into "quiz"(quiz_id, users, status, num_players, questions) values($1, $2, $3, $4, $5)`
	_, err := db.Exec(insertQuiz, quiz.Game_ID, "{"+strings.Join(quiz.Users, ",")+"}",
		quiz.Status, quiz.NumPlayers, "{"+strings.Join(transformQuestionsToString(quiz.Questions), ",")+"}")
	if err != nil {
		log.Printf("Errors %s insering Question to DB", err)
	}
	log.Printf("Insert quiz to DB successfully")
}

func deleteQuizFromDatabase(db *sql.DB, quiz_id string) {
	deleteQuestion := `delete from "quiz" where quiz_id=$1`
	_, err := db.Exec(deleteQuestion, quiz_id)

	if err != nil {
		log.Printf("Errors %s deleting Quiz to DB", err)
	}
	log.Printf("Deleted quiz with ID %s successfully", quiz_id)
}

func deleteQuizzesFromDatabase(db *sql.DB) {
	deleteQuestion := `delete from "quiz"`
	_, err := db.Exec(deleteQuestion)

	if err != nil {
		log.Printf("Errors %s deleting all Quiz game from DB", err)
	}
	log.Printf("Deleted all quiz game from DB successfully")
}