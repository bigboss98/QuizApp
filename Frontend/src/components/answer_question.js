import { AnswerQuestion } from "../api/start_game";

export default async function answerQuestion(websocket, userName, roomName, choice, timeAnswer){
    /* 
     * Component to answer Question given quiz_id and the answer choice
     *
     * Params:
     *      websocket: Websocket connection to communicate to the server  
     *      choice: Answer choice of a given question 
     * Return a JSON object with the following semantics
     * {
     *      "correct_answer": correct answer retrieved from server 
     *      "guess": a boolean that indicate whether response was guessed by user 
     * }
     */
    const answer = async () => {
        await AnswerQuestion(websocket, userName, roomName, {
                "User": userName,
                "Answers": [choice],
                "TimeToAnswer": timeAnswer,
            });
    }
    answer();
    websocket.onmessage = (e) => {
        // a message was received
        var message = JSON.parse(e.data)
        var received_message = JSON.parse(message.Message)    
        return received_message
    }
}