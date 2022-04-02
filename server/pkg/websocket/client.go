package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/riversy/chat-room/server/pkg/messages"
	"github.com/riversy/chat-room/server/pkg/users"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Uuid string
	Conn *websocket.Conn
	Pool *Pool
	User *users.User
}

func NewClient(conn *websocket.Conn, pool *Pool, user *users.User) *Client {
	return &Client{
		Uuid: uuid.NewString(),
		Conn: conn,
		Pool: pool,
		User: user,
	}
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, messageBody, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		transport := &MessageTransport{}
		err = json.Unmarshal(messageBody, transport)
		if err != nil {
			log.Println(err)
			continue
		}

		switch transport.MessageType {
		case messages.LoginRequestMessage:
			loginRequest, _ := transport.Payload.(*messages.LoginRequestMessageObject)
			err := c.Pool.BroadcastServerMessage(fmt.Sprintf("%s Connected", loginRequest.Username))
			if err != nil {
				fmt.Println(err)
				break
			}
			c.User = users.NewUser(loginRequest.Username)
			c.Pool.LogIn <- c
			fmt.Printf("LoginRequestMessageObject Received from: %s\n", c.User.Name)
			break
		case messages.LogoutRequestMessage:
			logoutRequest, _ := transport.Payload.(*messages.LogoutRequestMessageObject)
			userUuid, err := users.GetUserUuidFromJwt(logoutRequest.JWT)
			if err != nil {
				fmt.Printf("It's impossible extract user from JWT")
				continue
			}

			if userUuid != c.User.Uuid {
				fmt.Printf("Invalid request of logout")
				continue
			}

			c.Pool.LogOut <- c
			break
		case messages.ClientTextMessage:
			if c.User == nil {
				fmt.Printf("Message from unauthorized user: '%s'. Skipping...\n", string(messageBody))
				continue
			}
			clientMessage, _ := transport.Payload.(*messages.ClientTextMessageObject)
			payload := messages.NewServerTextMessage(c.User.Name, clientMessage.Text)
			transport := NewMessageTransport(messages.ServerTextMessage, payload)
			fmt.Println(c, c.User)
			fmt.Printf("ClientTextMessage Received: %s %+v\n", c.User.Name, messageBody)
			c.Pool.Broadcast <- *transport
			c.Pool.History.AddMessage(payload)
			break
		default:
			fmt.Println("Incoming message is unknown:", string(messageBody))
		}

	}
}

// IsAuthorized identify if the client provided its name before.
func (c *Client) IsAuthorized() bool {
	return c.User != nil
}
