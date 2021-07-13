/*
 * Api functions to start a Quiz game 
 */
import axios from 'axios';

export async function startQuiz(requestHeader) {
    /*
     * API function that make a request to start a Quiz game giving also the first question
     * 
     * Param:
     *      requestHeader: JSON object that contains params 
     */
    console.log(requestHeader)
    const response = await axios.post('http://192.168.1.75:8080/start_quiz', JSON.stringify(requestHeader));
    return response.data;
}

export async function getQuestion(gameId){
    /*
     * API function that make a request to obtain the next question giving a gameId quiz game
     */
    console.log(gameId)
    const response = await axios.get('http://192.168.1.75:8080/get_question/' + gameId);
    return response.data;
}

export async function AnswerQuestion(gameId, requestHeader){
    const response = await axios.post('http://192.168.1.75:8080/answer_question/' + gameId,
                                      JSON.stringify(requestHeader));
    return response.data; 
}