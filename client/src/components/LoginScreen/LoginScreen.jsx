import React, {Component} from "react";
import "./LoginScreen.scss";

class LoginScreen extends Component {
    renderParticipantsSection() {
        const participantsText = this.props.participantsQty === 0 ?
            `There is nobody in the room.` :
            `There are ${this.props.participantsQty} bodies in the room`

        return (
            <div className='participants-section'>
                <h3>{participantsText}</h3>
            </div>
        )
    }

    render() {
        return (
            <div className='login-screen'>
                <h2>Welcome!</h2>
                {this.renderParticipantsSection()}

                <div className="login-input">
                    <input placeholder={"Type your name and press Enter to log in..."}
                           onKeyDown={this.props.logIn.bind(this)} />
                </div>

            </div>
        );
    };
}

export default LoginScreen;