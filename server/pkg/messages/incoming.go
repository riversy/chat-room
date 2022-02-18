package messages

// LoginRequestMessageObject is sent by the client in the case to enter to the chat room.
type LoginRequestMessageObject struct {
	Username string `json:"username"`
}

// LogoutRequestMessageObject is sent by the server to the client.
// It contains JWT token with user.User identifier. It will allow identifying user.User
// if he will connect to the chat later or from other tab.
type LogoutRequestMessageObject struct {
	JWT string `json:"jwt"`
}

// ClientTextMessageObject contains one text message sent by the client.
type ClientTextMessageObject struct {
	Text string `json:"text"`
}
