import React, {useState, useEffect} from 'react';
import { StyleSheet, View, Text, TouchableOpacity, TextInput, Button } from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import { LoginRequest } from '../api/start_game';

export default function Login({ navigation }) {
    /*
     * SignUp is the component used to register an User of Quizzone App
     * 
     * Params:
     * -navigation: navigation object used to navigate between screens 
     */
    const [name, setName] = useState("")
    const [password, setPassword] = useState("")

    return (
        <View style={styles.container}>
            <View style={{flexDirection:"row"}}>
                <Text>Name: </Text>
                <TextInput onChangeText={setName} value={name}></TextInput>
            </View>
            <View style={{flexDirection:"row"}}>
                <Text>Password: </Text>
                <TextInput onChangeText={setPassword} value={password} secureTextEntry={true}></TextInput>
            </View>
            <Button title="Sign Up" onPress={() => {
                const login = async () => {
                    const data = await LoginRequest(name, password)
                    console.log(data)
                    navigation.navigate('Home', {
                        name: data.name,
                        token: data.token,
                    })
                }
                login()
                
            }}></Button>    
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