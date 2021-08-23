import React, {useState, useEffect} from 'react';
import { StyleSheet, View, Text, TouchableOpacity, TextInput, Button } from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import { SignUpRequest } from '../api/start_game';

export default function SignUp({ navigation }) {
    /*
     * SignUp is the component used to register an User of Quizzone App
     * 
     * Params:
     * -navigation: navigation object used to navigate between screens 
     */
    const [name, setName] = useState("")
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
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
            <View style={{flexDirection:"row"}}>
                <Text>Confirm Password: </Text>
                <TextInput onChangeText={setConfirmPassword} value={confirmPassword} secureTextEntry={true}></TextInput>
            </View>
            <Button title="Sign Up" onPress={() => {
                if (password === confirmPassword) {
                    const signUp = async () => {
                        return await SignUpRequest(name, password)
                    }
                    const data = signUp()
                    navigation.navigate('Home', {
                        token: "",
                        name: ""
                    })
                }
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