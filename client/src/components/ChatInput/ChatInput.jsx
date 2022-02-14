import React, { Component } from "react";
import "./ChatInput.scss";

class ChatInput extends Component {
    render() {
        return (
            <div className="chat-input">
                <input placeholder={"Type message and press Enter to send..."} onKeyDown={this.props.send} />
            </div>
        );
    }
}

export default ChatInput;