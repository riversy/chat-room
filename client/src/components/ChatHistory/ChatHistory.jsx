import React, {Component} from "react";
import Message from "../Message"
import "./ChatHistory.scss";

class ChatHistory extends Component {
    render() {
        const messages = this.props.chatHistory.map(msg => <Message key={msg.uuid} message={msg}/>);
        return (
            <div className='chat-history'>
                <h2>Chat History</h2>
                {messages}
            </div>
        );
    };
}

export default ChatHistory;