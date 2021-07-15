import React from 'react';
import { View, Text, Button } from 'react-native';
import { timeToAnswer } from '../common/common';
import { Choice, getStyleChoice } from '../components/question';

export default function AnswerQuestion({ route, navigation }) {
    /*
     * AnswerQuestion is the scene used to show Answer of a Question
     *
     * Param:
     */
        const {quizId, answer_given, choices, question, answer} = route.params;
        return (
            <View>
                <Text> {question}</Text>
                <View style={{flexDirection: 'row'}}>
                    <Choice styles={getStyleChoice(answer, answer_given, choices[0])} choice={choices[0]}
                            quizId={quizId} question={question} choices={choices}
                            navigation={navigation}></Choice>
                    <Choice styles={getStyleChoice(answer, answer_given, choices[1])} choice={choices[1]}
                            quizId={quizId} question={question} choices={choices}
                            navigation={navigation}></Choice>
                </View>
                <View style={{flexDirection: 'row'}}>
                    <Choice styles={getStyleChoice(answer, answer_given, choices[2])} choice={choices[2]} 
                            quizId={quizId} question={question} choices={choices}
                            navigation={navigation}></Choice>
                    <Choice styles={getStyleChoice(answer, answer_given, choices[3])} choice={choices[3]}
                            quizId={quizId} question={question} choices={choices}
                            navigation={navigation}></Choice>
                </View>
                <Button title="Next" onPress={()=>
                        navigation.navigate('QuizMatch', {
                            gameId: quizId,
                            started: false,
                            questionTimer: timeToAnswer,
                        })
                    }/>  
            </View>
        );
}

