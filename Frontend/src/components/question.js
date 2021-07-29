/*
 * Component to render a question for Quiz Game 
 */
import React, {useState, useEffect} from 'react';
import { StyleSheet, Text, View, TouchableOpacity, Button} from 'react-native';
import answerQuestion from './answer_question';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import { timeToAnswer } from '../common/common';

export default function Question({question, choices, navigation, quizId, seconds}){
    /*
     * Component to render a Question of the Quiz game 
     * Params:
     *      question: current Question text 
     *      choices: Array of possible answer choices
     *      navigation: navigation object to navigate to different Scenes
     *      quizId: ID of Quiz match
     */
    
    return (
        <View> 

            <Text> {question}</Text>
            <View style={{flexDirection: 'row'}}>
                <Choice style={styles.choice} choice={choices[0]} 
                        quizId={quizId} question={question} choices={choices}
                        navigation={navigation} seconds={seconds}></Choice>
                <Choice style={styles.choice} choice={choices[1]}
                        quizId={quizId} question={question} choices={choices}
                        navigation={navigation} seconds={seconds}></Choice>
            </View>
            <View style={{flexDirection: 'row'}}>
                <Choice style={styles.choice} choice={choices[2]}
                        quizId={quizId} question={question} choices={choices}
                        navigation={navigation} seconds={seconds}></Choice>
                <Choice style={styles.choice} choice={choices[3]}
                        quizId={quizId} question={question} choices={choices}
                        navigation={navigation} seconds={seconds}></Choice>
            </View>
        </View>
    );
}

export function Choice({style, choice, quizId, navigation, question, choices, seconds}){
    /*
     * Component to render a Choice of a Question 
     * Params:
     *      style: style to use to render a Choice
     *      choice: answer choice to represent
     *      quizId: ID of Quiz match
     *      navigation: navigation object to navigate in case 
     *      question: Question text used to pass to AnswerQuestion scene
     *      choices: array of Answer choices 
     */
    return (
        <TouchableOpacity styles={style}
                          onPress={() => {
                              navigation.navigate('AnswerQuestion', {
                                quizId: quizId,
                                answer_given: choice, 
                                choices: choices,
                                question: question,
                                answer: answerQuestion(quizId, choice, timeToAnswer - seconds)
                              })
                          }}>

            <Text> {choice}</Text>
        </TouchableOpacity>
    )
}

export function getStyleChoice(answer, answer_given, choice){
    /*
     * Function used to select Style to render Choice based on 
     * correct answer and answer chosen
     * 
     * Param:
     *      answer: JSON returned by answerQuestion API
     *      answer_given: answer given by a user
     *      choice: answer choice to choose the style 
     */
    if(answer_given != choice){
        return styles.choice;
    }
    if(answer_given == choice && answer.guess){
        return styles.correct;
    }
    return styles.wrong;
}

const styles = StyleSheet.create({
    correct: {
            //flex: 1,
            //flexDirection: 'row',
            //width: 100+"%",
            //height: 100+"%",
            borderRadius: wp('10%'),
            backgroundColor: 'green',
            //alignItems: 'center',
            //justifyContent: 'center',
    },
    wrong: {
            //flexDirection: 'row',
            borderRadius: wp('5%'),
            backgroundColor: 'red',
    },
    choice: {
        borderRadius: wp('5%'),
        backgroundColor: 'white',     
        //flex: 1,
        //resizeMode: "cover",
        //justifyContent: "center",
    },
    textTitle : {
        fontSize : wp('10%'),
        color: '#4280ff',
        fontWeight: 'bold'
    },
    buttonStyle : {
            backgroundColor: '#4280ff',
            //width: 30+"%",
            //height: 30+"%",
            //borderRadius : wp('10%'),
    },
    buttonText : {
        fontSize: wp('5%'),
    }
});