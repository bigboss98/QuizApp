import { AnswerQuestion } from "../api/start_game";

export default function answerQuestion(quiz_id, choice, timeAnswer){
    /* 
     * Component to answer Question given quiz_id and the answer choice
     *
     * Params:
     *      quiz_id: ID of Quiz game 
     *      choice: Answer choice of a given question 
     * Return a JSON object with the following semantics
     * {
     *      "correct_answer": correct answer retrieved from server 
     *      "guess": a boolean that indicate whether response was guessed by user 
     * }
     */
    const answer = async () => {
        const data = await AnswerQuestion(quiz_id, {
                "Answers": [choice],
                "TimeToAnswer": timeAnswer,
            });
        return {
            correct_answer: data.correct_answer,
            guess: data.guess, 
            score: data.score,            
        }
    }
    return answer();
}