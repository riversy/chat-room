package history

import "github.com/riversy/chat-room/server/pkg/messages"

type History struct {
	Messages []*messages.ServerTextMessageObject
}

func NewHistory(size int) *History {
	return &History{
		Messages: make([]*messages.ServerTextMessageObject, 0, size),
	}
}

func (h *History) AddMessage(message *messages.ServerTextMessageObject) {
	capacity := cap(h.Messages)
	if len(h.Messages) == capacity {
		for i := 0; i < capacity-1; i++ {
			h.Messages[i] = h.Messages[i+1]
		}
		h.Messages[capacity-1] = message
	} else {
		h.Messages = append(h.Messages, message)
	}
}

func (h *History) GetMessages() []*messages.ServerTextMessageObject {
	return h.Messages
}
