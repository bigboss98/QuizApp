import { AnswerQuestion } from "../api/start_game";

export default function answerQuestion(quiz_id, choice){

    const answer = async () => {
        const data = await AnswerQuestion(quiz_id, {
            "Answers": [choice],
            });
        return {
            correct_answer: data.correct_answer,
            guess: data.guess,            
            }
    }
    return answer();
}