import React, { Component } from "react";
import "./Message.scss";

class Message extends Component {
    constructor(props) {
        super(props);
        let messageObject = JSON.parse(this.props.message);
        this.state = {
            message: messageObject
        };
    }

    render() {
        return <div className="Message">{this.state.message.body}</div>;
    }
}

export default Message;