/*
 * Api functions to start a Quiz game 
 */
import axios from 'axios';
import { io } from "socket.io-client";

export function joinRoom(websocket, userName, roomName) {
    websocket.send(JSON.stringify({
        Target: {
            name: roomName
        },
        Sender: {
            name: userName
        },
        Action: "join-room"

    }))
}

export async function startSinglePlayer(websocket, userName, roomName){
    websocket.send(JSON.stringify({
        Target: {
            name: roomName
        },
        Sender: {
            name: userName
        },
        Action: "join-room"

    }))
    websocket.send(JSON.stringify({
        Target: {
            name: roomName
        },
        Sender: {
            name: userName
        },
        Action: "start-game"

    }))
}
export async function leaveRoom(websocket, userName, roomName) {

}
export async function startQuiz(websocket, userName, roomName) {
    /*
     * API function that make a request to start a Quiz game giving also the first question
     * 
     * Param:
     *      websocket: Websocket connection for a specific user
     *      username: Name of user that has websocket connection   
     */
    websocket.send(JSON.stringify({
        Target: {
            name: roomName
        },
        Sender: {
            name: userName
        },
        Action: "start-game"

    }))
}

export async function handleMessage(message) {

}
export async function getQuestion(websocket, userName, roomName){
    /*
     * API function that make a request to obtain the next question giving a gameId
     * 
     * Param:
     *      gameId: ID of Quiz Match to use to retrieve question 
     */
    const get_question = async () => {
        websocket.send(JSON.stringify({
            Target: {
                name: roomName
            },
            Sender: {
                name: userName
            },
            Action: "get-question"

        }))
    }
    get_question()
    
}

export async function AnswerQuestion(websocket, userName, roomName, requestBody){
    /*
     * API function to make request to answer Question giving gameId and a request Body with
     * answer of Question
     * 
     * Param:
     *      gameId: ID of Quiz Match to use to answer question
     *      requestBody: JSON object that contains answer of Question 
     */
    websocket.send(JSON.stringify({
        Target: {
            name: roomName
        },
        Sender: {
            name: userName
        },
        Message: requestBody,
        Action: "answer-question" 
    }))
}

export async function EndGame(websocket, userName, roomName){
    const end_game = async () => {
        websocket.send(JSON.stringify({
            Target: {
                name: roomName
            },
            Sender: {
                name: userName
            },
            Action: "end-game" 
        }))
    }
    end_game()

}

export async function GetPlayers(websocket, userName, roomName) {
    websocket.send(JSON.stringify({
       Target: {
           name: roomName
       },
       Sender: {
           name: userName
       },
       Action: "get-players"
    })) 
}

export async function SignUpRequest(username, password){
    const request = axios.post('http://192.168.1.75:8080/sign_up', JSON.stringify({
        name: username,
        password: password, 
        role: "User"
    }))   
    return request.data
}

export async function LoginRequest(username, password) {
    console.log("Dio")
    const request = await axios.post('http://192.168.1.75:8080/sign_in', JSON.stringify({
        name: username,
        password: password, 
        role: "User",
    }))
    console.log("DIo")
    let data = await request.data 
    console.log("AAA")
    return data 
}

export async function validateSocket(websocket, name, token) {
    const validate = async () => {

        await websocket.send(JSON.stringify({
            Action: "auth",
            Token: token,
            Sender: {
                name: name
            },
            Target: {
                name: ""
            }
        }))
    }

    validate()
    websocket.onmessage = (e) => {
        const status_authentification = e.data 
        if (status_authentification === "User Authorized") {
            return true
        }else{
            return false
        } 
    }

}