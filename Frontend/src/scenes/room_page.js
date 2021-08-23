import { View, Text, Button, FlatList} from "react-native";
import { GetPlayers, startQuiz } from "../api/start_game";
import React, {useState, useEffect} from "react";

export default function RoomPage({route, navigation}) {
    const {websocket, userName, roomName} = route.params 
    
    const [items, setItems] = useState([])

    const [timer, setTimer] = useState({
        seconds: 0,
    })

	
	useEffect(()=> {
		setTimer({
			seconds: 0,
		})
	}, [route])
    useEffect(() => {
        let timeout = setTimeout(() => {
            const getPlayers = async () => { 
                await GetPlayers(websocket, userName, roomName)
            }
            getPlayers()
            websocket.onmessage = (e) => {
                // a message was received
                console.log("DIO CANE")
                var message = JSON.parse(e.data)
                var received_message = JSON.parse(message.Message)
                console.log(received_message)
                setItems(received_message.players)
            }
         }, 1000);
        return () => clearTimeout(timeout);
    }, [items])
    return (
        <View>
                <FlatList
                    data={items}
                    renderItem={({ item }) => <Text>Player {item.Name}: {item.Status}</Text>}
                    keyExtractor={(item) => item.name}
                /> 
                <Button title="Start Game" onPress={() => {
                    const start_quiz = async () => {
                        await startQuiz(websocket, userName, roomName)
                    }
                    start_quiz()
                    websocket.onmessage = (e) => {
                        var message = JSON.parse(e.data)
                        var received_message = JSON.parse(message.Message)
                        var question = received_message.question.Question
                        var choices = received_message.question.Choices
                        navigation.navigate('QuizMatch', {
                            websocket: websocket, 
                            question: question,
                            choices: choices,
                            userName: userName,
                            roomName: roomName,
                            started: true,

                        })
                    }
                }}></Button>
        </View>
    );
}