package messages

type MessageType int

const (
	_ MessageType = iota
	LoginRequestMessage
	LoginResponseMessage
	LogoutRequestMessage
	LogoutResponseMessage
	UserEnterMessage
	UserQuitMessage
	ClientTextMessage
	ServerTextMessage
	ParticipantsQtyUpdate
)
