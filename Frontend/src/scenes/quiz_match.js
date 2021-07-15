import React from 'react';
import { StyleSheet, Text, View, TouchableOpacity, Button} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import { useState, useEffect } from 'react';
import { getQuestion } from '../api/start_game';
import Question from '../components/question.js';
import { State } from 'react-native-gesture-handler';
import {timeToAnswer} from '../common/common.js'

export default function QuizMatch({ route, navigation }) {
	/*
	 * HomeScreen is the component used to render the home of Briscola game App
	 */

		const {gameId, started, questionTimer} = route.params;

		const [quiz, setQuiz] = useState({
			started: started,
			quiz_id: gameId,
			question: {},
		});

		let resetTimer = useEffect(() => {
				setTimer({
					seconds: timeToAnswer,
				})
			}, [route])

		useEffect(() => {
			const getMyQuestion = async () => {
				const data = await getQuestion(gameId)
				if(data.status){
					navigation.navigate('EndQuiz')
				}
				setQuiz({ 
					started: true,
					quiz_id: gameId,
					question: data.question,
				});
	
			}
			getMyQuestion()
		}, [route])

		const [timer, setTimer] = useState({
            seconds: questionTimer,
        })

	
	useEffect(()=> {
		setTimer({
			seconds: timeToAnswer,
		})
	}, [])
    useEffect(() => {
        let timeout = setTimeout(() => {
            console.log(timer.seconds)
            if(timer.seconds > 0){
                setTimer({
                    seconds: timer.seconds - 1,
                })
            }else{
                setTimer({
                    seconds: 40,
                })
                navigation.navigate('AnswerQuestion', {
                    quizId: gameId,
                    answer_given: "", 
                    choices: quiz.question.Choices,
                    question: quiz.question.Question,
                    answer: answerQuestion(gameId, "")
                })
            }
         }, 1000);
        return () => clearTimeout(timeout);
    })
	return (
		<View style={styles.container}>
			{ quiz.started ?             
					<View>
						<Text>Remaining Time: {timer.seconds} </Text>
						<Question question = {quiz.question.Question} choices={quiz.question.Choices}
							  	  navigation={navigation} quizId={quiz.quiz_id} resetTimer={resetTimer}></Question>
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