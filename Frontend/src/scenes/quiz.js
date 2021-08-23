import 'react-native-gesture-handler';
import * as React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import createStackNavigator from "@react-navigation/stack/src/navigators/createStackNavigator";

import HomeScreen from './homepage.js';
import StartGame from './game.js';
import InsertQuestion from './insert_question.js';
import QuizMatch from './quiz_match.js';
import AnswerQuestion from './answer_question.js';
import EndQuiz from './end_quiz.js';
import RoomPage from './room_page.js';
import Login from './login.js';
import SignUp from './register.js';

const Stack = createStackNavigator();

export default function Quiz(){
    /*
     * Quiz function is a scenes function used to create the Quiz game app
     * It contains 4 different scenes:
     * -HomeScreen: scene to render the home of the game
     * -StartGame: scene used to render start of a Quiz game 
     * -InsertQuestion: scene used to propose custom question and answer  
     * -QuizMatch: scene used to play a Quiz game 
     * -AnswerQuestion: scene used to answer Question of a Quiz game
     * -EndQuiz: scene used to render end of a Quiz game  
     * 
     *                 
     */
    return (
        <NavigationContainer>
            <Stack.Navigator>
                <Stack.Screen
                    name="Home"
                    component={HomeScreen}
                    initialParams={{token: "", name: ""}}
                />
                <Stack.Screen
                    name="Sign Up"
                    component={SignUp}
                />
                <Stack.Screen 
                    name="Login"
                    component={Login}
                />
                <Stack.Screen
                    name="StartGame"
                    component={StartGame}
                />
                <Stack.Screen
                    name="InsertQuestion"
                    component={InsertQuestion}
                />
                <Stack.Screen
                    name="RoomPage"
                    component={RoomPage}
                />
                <Stack.Screen
                    name="QuizMatch"
                    component={QuizMatch}
                />
                <Stack.Screen 
                    name="AnswerQuestion"
                    component={AnswerQuestion}
                />
                <Stack.Screen
                    name="EndQuiz"
                    component={EndQuiz}
                />
            </Stack.Navigator>
        </NavigationContainer>
    )
}
