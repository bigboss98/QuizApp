import { StatusBar } from 'expo-status-bar';
import React, {useState, useEffect} from 'react';
import { StyleSheet, Text, View, Button, TouchableOpacity} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import { startQuiz } from '../api/start_game';
import { timeToAnswer } from '../common/common';

export default function StartGame({ navigation }) {
    /*
     * SelectGameType is the component used to select which modalities of Briscola game between
     * Single player and Multiplayer
     */
        return (
            <View style={styles.container}>
                <TouchableOpacity style={styles.buttonStyle}
                    onPress={() => {
                            console.log("Dio")
                            const startMyGame = async () => {
                                const data = await startQuiz({
                                    users: ["Marco"],
                                })
                                navigation.navigate('QuizMatch', {
                                    gameId: data.game_id,
                                    started: false,
                                    questionTimer: timeToAnswer,
                                })
                            }
                            startMyGame()
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