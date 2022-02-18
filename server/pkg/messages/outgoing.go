package messages

import "github.com/google/uuid"

// LogoutResponseMessageObject is sent by the server to the client.
// It's an empty object which indicates that the logout is acknowledged.
type LogoutResponseMessageObject struct {
}

func NewLogoutResponseMessageObject() *LogoutResponseMessageObject {
	return &LogoutResponseMessageObject{}
}

// LoginResponseMessageObject is sent by the server to the client.
// It contains JWT token with user.User identifier. It will allow identifying user.User
// if he will connect to the chat later or from other tab.
type LoginResponseMessageObject struct {
	JWT string `json:"jwt"`
}

func NewLoginResponseMessageObject(jwt string) *LoginResponseMessageObject {
	return &LoginResponseMessageObject{JWT: jwt}
}

// ServerTextMessageObject contains one text message to build the conversation
// sent by the server.
type ServerTextMessageObject struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

// NewServerTextMessage generates new MessageText Object to be sent
func NewServerTextMessage(username string, message string) *ServerTextMessageObject {
	return &ServerTextMessageObject{
		Uuid:     uuid.NewString(),
		Username: username,
		Text:     message,
	}
}

// ParticipantsQtyUpdateObject contains data about participants qty.
type ParticipantsQtyUpdateObject struct {
	Qty int `json:"qty"`
}

func NewParticipantsQtyUpdateObject(qty int) *ParticipantsQtyUpdateObject {
	return &ParticipantsQtyUpdateObject{Qty: qty}
}
