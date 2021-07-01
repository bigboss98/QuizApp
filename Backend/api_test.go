package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/start_game", startGame).Methods("POST")
	return router
}

func TestStartGame(test *testing.T) {
	users_json := []byte(`{
		"users": [
			{
				"name": "Marco"
			},
			{
				"name": "William"
			}
		]
	}`)
	request, _ := http.NewRequest("POST", "/start_game", bytes.NewBuffer(users_json))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(test, 200, response.Code, "OK response is expected")
}

/*
 *
 	func TestInsertQuestion(test *testing.T) {

}*/
