import React, {useState, useEffect} from 'react';
import { StyleSheet, Text, View, TouchableOpacity, Button} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import {FinalResults} from '../api/start_game';

export default function EndQuiz({ route, navigation}){
    const {gameId} = route.params
    const [finalResults, setFinalResults] = useState({
            status: "",
            winner: "",
            score: 0,
        }) 
	useEffect(() => {
		const finalMyResults = async () => {
			const data = await FinalResults(gameId);
			setFinalResults({ 
				status: data.status,
				winner: data.winner,
				score: data.score,
		    });
	
		}
		finalMyResults()
	}, [route])    
    return (
            <View>
                {finalResults.status != "" ?
                    <View>
                        <Text> End Quiz </Text> 
                        <Text> Score: {finalResults.score} </Text>
                        <Button title="goBack" onPress={() => {
                            navigation.navigate('Home')
                            }}></Button>
                    </View> : <Text></Text>}
             </View>
        );
}


