import React from 'react';
import { View, Text, Button } from 'react-native';
import { getQuestion } from '../api/start_game';
import { timeToAnswer } from '../common/common';
import { Choice, getStyleChoice } from '../components/question';

export default function AnswerQuestion({ route, navigation }) {
    /*
     * AnswerQuestion is the scene used to show Answer of a Question
     *
     * Param:
     */
        const {websocket, answer_given, choices, question, answer} = route.params;
        return (
            <View>
                <Text> {question}</Text>
                <View style={{flexDirection: 'row'}}>
                    <Choice styles={getStyleChoice(answer, answer_given, choices[0])} choice={choices[0]}
                            choices={choices} navigation={navigation}></Choice>
                    <Choice styles={getStyleChoice(answer, answer_given, choices[1])} choice={choices[1]}
                            question={question} choices={choices} navigation={navigation}></Choice>
                </View>
                <View style={{flexDirection: 'row'}}>
                    <Choice styles={getStyleChoice(answer, answer_given, choices[2])} choice={choices[2]} 
                            question={question} choices={choices} navigation={navigation}></Choice>
                    <Choice styles={getStyleChoice(answer, answer_given, choices[3])} choice={choices[3]}
                            question={question} choices={choices} navigation={navigation}></Choice>
                </View>
                <Button title="Next" onPress={()=> {
                    var get_question = async () => {
                        await getQuestion(websocket, "Marco", "room1")
                    }
                    get_question()
                    websocket.onmessage = (e) => {
                        // a message was received
                        var message = JSON.parse(e.data)
                        var received_message = JSON.parse(message.Message)
                        console.log(received_message.question)
                        var question = received_message.question 
                        
                        navigation.navigate('QuizMatch', {
                            websocket: websocket, 
                            question: question.Question,
                            choices: question.Choices,
                            userName: "Marco",
                            roomName: "room1",
                            started: true,
                    })}}}></Button>
            </View>
        );
}

