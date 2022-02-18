import React, { Component } from "react";
import "./Message.scss";

class Message extends Component {
    constructor(props) {
        super(props);
        this.state = {
            message: props.message
        };
    }

    render() {
        return <div key={this.state.message.uuid} className="Message">{this.state.message.username}: {this.state.message.text}</div>;
    }
}

export default Message;