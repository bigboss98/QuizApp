import React from 'react';
import { StyleSheet, Text, View, TouchableOpacity, Button} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import { useState, useEffect } from 'react';
import { startQuiz, getQuestion } from '../api/start_game';
import Question from '../components/question.js';
import { State } from 'react-native-gesture-handler';

export default function QuizMatch({ route, navigation }) {
		/*
		 * HomeScreen is the component used to render the home of Briscola game App
		 */

		const {gameId, started} = route.params;

		const [quiz, setQuiz] = useState({
			started: started,
			quiz_id: gameId,
			question: {},
		});
		
		useEffect(() => {
			console.log("Dio")
			const getMyQuestion = async () => {
				const data = await getQuestion(gameId)
				console.log(data)
				console.log(status)
				if(data.status){
					console.log("END")
					navigation.navigate('EndQuiz')
					console.log("DIO")
				}
				setQuiz({ 
					started: true,
					quiz_id: gameId,
					question: data.question,
				});
	
			}
			getMyQuestion()
		}, [route])


		return (
				<View style={styles.container}>
					{ quiz.started ? 
						<Question question = {quiz.question.Question} choices={quiz.question.Choices}
								  navigation={navigation} quizId={quiz.quiz_id}></Question> :
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