import React, {Component} from "react";
import './App.scss';
import Header from './components/Header';
import ChatHistory from './components/ChatHistory';
import LoginScreen from './components/LoginScreen';
import ChatInput from './components/ChatInput';
import * as client from "./api/websocket/client";
import {LoginRequestMessage, LoginResponseMessage} from "./api/websocket/client";

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            userJWT: null,
            participantsQty: 0,
            chatHistory: [],
        };

        this.sendMessage = this.sendMessage.bind(this);
    }

    onConnected() {
        console.log("onConnected");
    }

    onDisconnected() {
        console.log("onDisconnected");
    }

    onError(error) {
        console.error(error);
    }

    onUserEnter(payload) {
        console.log("onUserEnter", {payload});
    }

    onUserQuit(payload) {
        console.log("onUserQuit", {payload});
    }

    onPeopleUpdate(payload) {
        console.log("onPeopleUpdate", {payload});
    }

    onLogout() {
        this.setState({
            userJWT: null
        })
    }

    onLogin(payload) {
        const {jwt = null} = payload;
        this.setState({
            userJWT: jwt
        })
    }

    logIn(event) {
        if (event.keyCode === 13) {
            console.log(this);
            this.client.sendMessage(
                LoginRequestMessage,
                {"username": event.target.value}
            );
            event.target.value = "";
        }
    }

    onMessageHandler(payload) {
        console.log("New Message", {payload});
        this.setState(prevState => ({
            chatHistory: [...this.state.chatHistory, payload]
        }));
    };

    onParticipantsQtyUpdate(payload) {
        const {qty} = payload;
        this.setState({
            participantsQty: qty
        });
    }

    componentDidMount() {
        let handlers = {};
        handlers[client.LoginResponseMessage] = this.onLogin.bind(this);
        handlers[client.LogoutResponseMessage] = this.onLogout.bind(this);
        handlers[client.UserEnterMessage] = this.onUserEnter.bind(this);
        handlers[client.UserQuitMessage] = this.onUserQuit.bind(this);
        handlers[client.ServerTextMessage] = this.onMessageHandler.bind(this);
        handlers[client.ParticipantsQtyUpdate] = this.onParticipantsQtyUpdate.bind(this);
        handlers[client.ConnectedEvent] = this.onConnected.bind(this);
        handlers[client.DisconnectedEvent] = this.onDisconnected.bind(this);
        handlers[client.ErrorEvent] = this.onError.bind(this);

        this.client = new client.WebsocketClient(
            process.env.REACT_APP_WS_PATH || `ws://localhost:8080/ws`,
            handlers
        );
    }

    sendMessage(event) {
        if (event.keyCode === 13) {
            this.client.sendMessage(client.ClientTextMessage, {"text": event.target.value});
            event.target.value = "";
        }
    }

    renderLoginScreen() {
        return (
            <LoginScreen participantsQty={this.state.participantsQty} logIn={this.logIn.bind(this)}  />
        )
    }

    renderChatHistory() {
        return (
            <div>
                <ChatHistory chatHistory={this.state.chatHistory}/>
                <ChatInput send={this.sendMessage}/>
            </div>
        )
    }

    render() {
        const ChatMainElement = this.state.userJWT ?
            this.renderChatHistory() :
            this.renderLoginScreen()

        return (
            <div className="App">
                <Header/>
                {ChatMainElement}
            </div>
        );
    }
}

export default App;
