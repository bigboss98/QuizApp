/*
 * Api functions to start a Quiz game 
 */
import axios from 'axios';

export async function startQuiz(requestBody) {
    /*
     * API function that make a request to start a Quiz game giving also the first question
     * 
     * Param:
     *      requestBody: JSON object that contains params 
     */
    const response = await axios.post('http://192.168.1.75:8080/start_quiz',
                                      JSON.stringify(requestBody));
    return response.data;
}

export async function getQuestion(gameId){
    /*
     * API function that make a request to obtain the next question giving a gameId
     * 
     * Param:
     *      gameId: ID of Quiz Match to use to retrieve question 
     */
    const response = await axios.get('http://192.168.1.75:8080/get_question/' + gameId);
    return response.data;
}

export async function AnswerQuestion(gameId, requestBody){
    /*
     * API function to make request to answer Question giving gameId and a request Body with
     * answer of Question
     * 
     * Param:
     *      gameId: ID of Quiz Match to use to answer question
     *      requestBody: JSON object that contains answer of Question 
     */
    const response = await axios.post('http://192.168.1.75:8080/answer_question/' + gameId,
                                      JSON.stringify(requestBody));
    return response.data; 
}