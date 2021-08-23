import { StatusBar } from 'expo-status-bar';
import React, {useState, useEffect} from 'react';
import { StyleSheet, Text, View, Button, TouchableOpacity} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import { startQuiz, startSinglePlayer } from '../api/start_game';
import { timeToAnswer } from '../common/common';

export default function StartGame({ route, navigation }) {
    /*
     * SelectGameType is the component used to select which modalities of Briscola game between
     * Single player and Multiplayer
     */
        const {name, token, websocket} = route.params

        return (
            <View style={styles.container}>
                <TouchableOpacity style={styles.buttonStyle}
                    onPress={() => {
                            console.log("Dio")
                            const startMyGame = async () => {
                                await startSinglePlayer(websocket, name, "room1")
 
                            }
                            startMyGame()
                            websocket.onmessage = (e) => {
                                // a message was received
                                var message = JSON.parse(e.data)
                                var received_message = JSON.parse(message.Message)
                                var question = received_message.question.Question
                                var choices = received_message.question.Choices
                                navigation.navigate('QuizMatch', {
                                    websocket: websocket, 
                                    question: question,
                                    choices: choices,
                                    userName: name,
                                    roomName: "room1",
                                    started: true,

                                })
                            }
                    }}>
                    <Text style={styles.buttomText}>Start Game</Text>
                </TouchableOpacity>
                
                <Button title="Go back" style={styles.buttomText} onPress={() => navigation.goBack()} />
            </View>
        );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        alignItems: 'center',
        marginTop: hp('15%'),
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
        marginTop: hp('5%'),
        backgroundColor: '#4280ff',
        padding: wp('10%'),
        borderRadius : wp('10%'),
    },
    buttonStyle2 : {
        marginTop: hp('5%'),
        backgroundColor: '#4280ff',
        padding: wp('10%'),
        borderRadius : wp('10%'),
    },
    buttonText : {
        fontSize: wp('10%'),
    }
});