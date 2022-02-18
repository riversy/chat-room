package websocket

import (
	"reflect"
	"testing"

	"github.com/riversy/chat-room/server/pkg/messages"
)

func TestMessageTransport_MarshalJSON(t *testing.T) {
	type fields struct {
		MessageType messages.MessageType
		Payload     interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"test ServerTextMessage marshaling",
			fields{
				messages.ServerTextMessage,
				&messages.ServerTextMessageObject{
					Uuid:     "test",
					Username: "Test",
					Text:     "Hello!",
				},
			},
			[]byte(`{"type":8,"payload":{"uuid":"test","username":"Test","text":"Hello!"}}`),
			false,
		},
		{
			"test LoginResponseMessage marshaling",
			fields{
				messages.LoginResponseMessage,
				&messages.LoginResponseMessageObject{
					JWT: "test",
				},
			},
			[]byte(`{"type":2,"payload":{"jwt":"test"}}`),
			false,
		},
		{
			"test ServerTextMessage marshaling",
			fields{
				messages.LogoutResponseMessage,
				&messages.LogoutResponseMessageObject{},
			},
			[]byte(`{"type":4,"payload":{}}`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MessageTransport{
				MessageType: tt.fields.MessageType,
				Payload:     tt.fields.Payload,
			}
			got, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestMessageTransport_UnmarshalJSON(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *MessageTransport
		wantErr bool
	}{
		{
			"test unmarshalling of LoginRequestMessage",
			args{
				[]byte(`{"type":1,"payload":{"username":"Test"}}`),
			},
			&MessageTransport{
				MessageType: messages.LoginRequestMessage,
				Payload: &messages.LoginRequestMessageObject{
					Username: "Test",
				},
			},
			false,
		},
		{
			"test unmarshalling of LogoutRequestMessage",
			args{
				[]byte(`{"type":3,"payload":{"jwt":"Test"}}`),
			},
			&MessageTransport{
				MessageType: messages.LogoutRequestMessage,
				Payload: &messages.LogoutRequestMessageObject{
					JWT: "Test",
				},
			},
			false,
		},
		{
			"test unmarshalling of ClientTextMessage",
			args{
				[]byte(`{"type":7,"payload":{"text":"Hello!"}}`),
			},
			&MessageTransport{
				MessageType: messages.ClientTextMessage,
				Payload: &messages.ClientTextMessageObject{
					Text: "Hello!",
				},
			},
			false,
		},
		{
			"test unmarshalling of CustomRequestMessage",
			args{
				[]byte(`{"type":50,"payload":{"field":"Test"}}`),
			},
			&MessageTransport{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MessageTransport{}

			err := m.UnmarshalJSON(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if (err != nil) && tt.wantErr {
				return
			}

			if !reflect.DeepEqual(m, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", m, tt.want)
			}
		})
	}
}

func TestNewMessageTransport(t *testing.T) {
	type args struct {
		messageType   messages.MessageType
		payloadObject interface{}
	}

	tests := []struct {
		name string
		args args
		want *MessageTransport
	}{
		{
			"test of new transport for Text Message creation",
			args{
				messageType:   messages.LogoutResponseMessage,
				payloadObject: messages.LogoutResponseMessageObject{},
			},
			&MessageTransport{
				messages.LogoutResponseMessage,
				messages.LogoutResponseMessageObject{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMessageTransport(tt.args.messageType, tt.args.payloadObject); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMessageTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}
