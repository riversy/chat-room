import React, { Component } from "react";
import "./ChatInput.scss";

class ChatInput extends Component {
    render() {
        return (
            <div className="chat-input">
                <input onKeyDown={this.props.send} />
            </div>
        );
    }
}

export default ChatInput;