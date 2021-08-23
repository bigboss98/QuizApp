import React from 'react';
import { StyleSheet, Text, View, TouchableOpacity, Button} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import { useState, useEffect } from 'react';
import { getQuestion, EndGame } from '../api/start_game';
import Question from '../components/question.js';
import { State } from 'react-native-gesture-handler';
import {timeToAnswer} from '../common/common.js'
import answerQuestion from '../components/answer_question';

export default function QuizMatch({ route, navigation }) {
	/*
	 * HomeScreen is the component used to render the home of Briscola game App
	 */

	const {websocket, question, choices, userName, roomName, started} = route.params;
	const [quiz, setQuiz] = useState({});

	useEffect(() => {
		setQuiz({
			started:started,
			choices: choices,
			question: question,
		})
		console.log(quiz)
		if(question == ""){
			var end_game = async () => {
				await EndGame(websocket, userName, roomName)
			}
			end_game()
			websocket.onmessage = (e) => {
				// a message was received
				var message = JSON.parse(e.data)
				var final_result = JSON.parse(message.Message)    
				
				console.log(final_result)
				navigation.navigate('EndQuiz', {
					websocket: websocket,
					results: final_result,
				})
				}}
	}, [route])

	const [timer, setTimer] = useState({
        seconds: timeToAnswer,
    })

	
	useEffect(()=> {
		setTimer({
			seconds: timeToAnswer,
		})
	}, [route])
    useEffect(() => {
        let timeout = setTimeout(() => {
            if(timer.seconds > 0){
                setTimer({
                    seconds: timer.seconds - 1,
                })
            }else{
                navigation.navigate('AnswerQuestion', {
					websocket: websocket,
                    answer_given: "", 
                    choices: quiz.choices,
                    question: quiz.question,
                    answer: answerQuestion(websocket, userName, roomName, "", timer.seconds)
                })
            }
         }, 1000);
        return () => clearTimeout(timeout);
    })
	return (
		<View style={styles.container}>
			{ quiz.started && quiz.question != "" ?             
					<View>
						<Text>Remaining Time: {timer.seconds} </Text>
						<Question question = {quiz.question} choices={quiz.choices}
							  	  navigation={navigation} seconds={timer.seconds}
								  websocket={websocket} username={userName}
								  roomName={roomName}></Question>
					</View> :
					<Text></Text>
			}
					
			 	</View>
			);
}

const styles = StyleSheet.create({
		container: {
				//flex: 1,
				flexDirection: 'row',
				//width: 100+"%",
				//height: 100+"%",
				//backgroundColor: '#1E4A62',
				//alignItems: 'center',
				//justifyContent: 'center',
		},
		buttomContainer: {
				flexDirection: 'row',
				//backgroundColor: '#1E4A62',
		},
		image: {
			flex: 1,
			resizeMode: "cover",
			justifyContent: "center",
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