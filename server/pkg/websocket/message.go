package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/riversy/chat-room/server/pkg/messages"
	"github.com/tidwall/gjson"
	"reflect"
)

type MessageMetadata struct {
	MessageType messages.MessageType `json:"type"`
}

type MessageTransport struct {
	MessageType messages.MessageType `json:"type"`
	Payload     interface{}          `json:"payload"`
}

func (m *MessageTransport) MarshalJSON() ([]byte, error) {
	payloadBytes, err := json.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}

	jsonString := fmt.Sprintf(
		`{"type":%d,"payload":%s}`,
		int(m.MessageType),
		string(payloadBytes),
	)

	return []byte(jsonString), nil
}

func (m *MessageTransport) UnmarshalJSON(bytes []byte) error {
	metadata, err := unmarshalMetadata(bytes)
	if err != nil {
		return err
	}

	m.MessageType = metadata.MessageType
	payload, err := unmarshalPayload(metadata.MessageType, bytes)
	if err != nil {
		return err
	}

	m.Payload = payload

	return nil
}

func unmarshalMetadata(bytes []byte) (*MessageMetadata, error) {
	result := &MessageMetadata{}
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func unmarshalPayload(messageType messages.MessageType, bytes []byte) (interface{}, error) {
	var result interface{}

	payloadBytes := gjson.GetBytes(bytes, "payload")
	if payloadBytes.Exists() {
		payloadType, err := getPayloadTypeFromMessageType(messageType)
		if err != nil {
			return nil, err
		}
		result = reflect.New(payloadType).Interface().(interface{})
		if err := json.Unmarshal([]byte(payloadBytes.Raw), &result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func getPayloadTypeFromMessageType(messageType messages.MessageType) (reflect.Type, error) {
	switch messageType {
	case messages.LoginRequestMessage:
		return reflect.TypeOf(messages.LoginRequestMessageObject{}), nil
	case messages.LogoutRequestMessage:
		return reflect.TypeOf(messages.LogoutRequestMessageObject{}), nil
	case messages.ClientTextMessage:
		return reflect.TypeOf(messages.ClientTextMessageObject{}), nil
	default:
		return nil, errors.New("message type is not defined")
	}
}

// NewMessageTransport generates new transport for the message.
func NewMessageTransport(messageType messages.MessageType, payloadObject interface{}) *MessageTransport {
	return &MessageTransport{
		MessageType: messageType,
		Payload:     payloadObject,
	}
}
