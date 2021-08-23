import React, {useState, useEffect} from 'react';
import { StyleSheet, Text, View, TouchableOpacity, Button} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import {FinalResults} from '../api/start_game';

export default function EndQuiz({ route, navigation}){
    const {websocket, results} = route.params
    const [finalResults, setFinalResults] = useState({
        status: results.status,
        winner: results.winner,
        scores: results.scores,
        players: results.players,
            
        }) 
    useEffect(() => {

        setFinalResults({
            status: results.status,
            winner: results.winner,
            scores: results.scores,
            players: results.players,
        })
	}, [route])   
    
    return (
            <View>
                {finalResults.status != "" ?
                    <View>
                        <Text> End Quiz </Text> 
                        <Text> Score {finalResults.players[0]}: {finalResults.scores[0]} </Text>
                        <Button title="goBack" onPress={() => {
                            navigation.navigate('Home')
                            }}></Button>
                    </View> : <Text></Text>}
             </View>
        );
}


