Endpoint implemented and available in QuizApp backend 

POST start_game()
    Start and initialize a new Quiz Game with params provided by a Request JSON with the following interface 
    {
        'users': [
            {
                'name': Name of users
            },
            {
                'name': Name of users
            }
            and so on having 2 or more users 
        ],
        'num_question': integer that represent the number of question of the quiz 
        'num_abstention': number of abstention given for each users at the start of quiz 

    }
    Return a JSON file formatted as follows 
    {
        'game_id': an UUID used to recognize uniquitely a QUiz game
        'num_question': number of question of the Quiz (> 0)
        'num_abstention': number of abstention available during the Quiz 
    }

GET get_question(string game_id)
    Get Question of Game_ID if available otherwise return empty string and 
    set status of game to "ended"

    JSON file used for respond as the following structure 
    {
        'status': indicates whether the QUiz is ended 
        'question': new Question of the game 
        'answer_choices': choices to answer the new question  
        'num_remaining_question': Number of remaining question 
    }

POST answer_question(string answer, string game_id)
    Answer to a question of a Quiz Game 
    Return with also the new question to ask to user if there are any otherwise 
    return a status that indicates the end of the game 

    JSON file returned by answer_question is the following 
    {  
        'status': indicates whether the Quiz is ended 
        'score': Score achieved up to now by User 
        'correct_answer': correct answer of the question 
    }
