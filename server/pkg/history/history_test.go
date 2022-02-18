package history

import (
	"os"
	"testing"

	"github.com/riversy/chat-room/server/pkg/messages"
	"github.com/stretchr/testify/assert"
)

var message1, message2, message3, message4 *messages.ServerTextMessageObject

func TestMain(m *testing.M) {
	message1 = messages.NewServerTextMessage("User 1", "Hello!")
	message2 = messages.NewServerTextMessage("User 2", "Hello!")
	message3 = messages.NewServerTextMessage("User 1", "How are you?")
	message4 = messages.NewServerTextMessage("User 2", "I'm fine!")

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestHistory_NewHistory(t *testing.T) {
	history := NewHistory(3)
	assert.Equal(t, 3, cap(history.Messages), "capacity expected to be 3")
}

func TestHistory_PushMessage_OneMessage(t *testing.T) {
	history := NewHistory(3)
	history.AddMessage(message1)

	assert.Equal(
		t,
		[]*messages.ServerTextMessageObject{message1},
		history.GetMessages(),
		"message1 expected to be only message in the history",
	)
}

func TestHistory_PushMessage_ManyMessages(t *testing.T) {
	history := NewHistory(3)
	history.AddMessage(message1)
	history.AddMessage(message2)
	history.AddMessage(message3)
	history.AddMessage(message4)

	assert.Equal(
		t,
		[]*messages.ServerTextMessageObject{message2, message3, message4},
		history.GetMessages(),
		"message2, message3, message4 expected to be in the history",
	)
}
