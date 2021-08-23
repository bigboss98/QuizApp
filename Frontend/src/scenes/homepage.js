import React from 'react';
import { StyleSheet, Text, View, TouchableOpacity, Button} from 'react-native';
import {widthPercentageToDP as wp, heightPercentageToDP as hp} from 'react-native-responsive-screen';
import { validateSocket, joinRoom } from '../api/start_game';

export default function HomeScreen({ route, navigation }) {
    /*
     * HomeScreen is the component used to render the home of Briscola game App
     */
    const {token, name } = route.params 
    var websocket
    if (name != "" && token != "") {
        websocket = new WebSocket('ws://192.168.1.75:8080/start_quiz/' + name)

        websocket.onopen = (e) => {
            validateSocket(websocket, name, token)
        }
        console.log(websocket)
    }
    return (
            <View styles={styles.container}>
                <View style={{flexDirection: "row"}}>
                    <Text>{name}</Text>
                    <Button title="Sign Up" onPress={()=> navigation.navigate('Sign Up')}></Button>
                    <Button title="Sign In" onPress={()=> navigation.navigate('Login')}></Button>
                </View>
                <Text style={styles.textTitle}> {navigation.name} </Text>
                <View styles={{flexDirection:'row'}}>
                    <Button //style={styles.buttonStyle}
                            title="Start Single Player Game"
                            onPress={() =>
                                          navigation.navigate('StartGame', {
                                              name: name,
                                              token: token, 
                                              websocket: websocket
                                          })
                                      }
                    >
                    </Button>
                    <Button //style={styles.buttonStyle}
                            title="Start MultiPlayer game"
                            onPress={() => {
                                        joinRoom(websocket, "Marco", "room1")
                                        navigation.navigate('RoomPage', {
                                              websocket: websocket,
                                              userName: name,
                                              roomName: "room1",
                                              token: token,
                                        })
                                    }}
                    >
                    </Button>
                    <Button
                        title="Insert Question"
                        onPress={
                            () => navigation.navigate('InsertQuestion')
                        }
                        //style={styles.buttonStyle}
                    >
                    </Button>
                </View>
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