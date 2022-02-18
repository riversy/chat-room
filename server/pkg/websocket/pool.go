package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/riversy/chat-room/server/pkg/history"
	"github.com/riversy/chat-room/server/pkg/messages"
	"github.com/riversy/chat-room/server/pkg/users"
)

type Pool struct {
	History    *history.History
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	LogIn      chan *Client
	LogOut     chan *Client
	Broadcast  chan MessageTransport
}

func NewPool() *Pool {
	return &Pool{
		History:    history.NewHistory(100),
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		LogIn:      make(chan *Client),
		LogOut:     make(chan *Client),
		Broadcast:  make(chan MessageTransport),
	}
}

func (pool *Pool) NewClient(user *users.User, conn *websocket.Conn) *Client {
	return NewClient(conn, pool, user)
}

// ServeWs creates a new Client in the pool and starts reading of incoming history
func (pool *Pool) ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	uuid, err := getUuidFromRequest(r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
		return
	}
	user, err := users.GetUserByUuid(uuid)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	conn, err := UpgradeConnection(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := pool.NewClient(user, conn)

	pool.Register <- client
	if client.IsAuthorized() {
		pool.SendHistoryToClient(client)
	}

	client.Read()
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.AddClientToPool(client)
			pool.SendUsersQtyToClient(client)
			break
		case client := <-pool.Unregister:
			pool.DeleteClientFromPool(client)
			break
		case client := <-pool.LogIn:
			err := pool.LogInClient(client)
			if err != nil {
				fmt.Println(err)
				break
			}

			err = pool.BroadcastUsersQty()
			if err != nil {
				fmt.Println(err)
				break
			}

			err = pool.SendHistoryToClient(client)
			if err != nil {
				fmt.Println(err)
				break
			}
			break
		case client := <-pool.LogOut:
			pool.LogOutClient(client)
			pool.BroadcastUsersQty()
			break
		case broadcast := <-pool.Broadcast:
			pool.BroadcastMessage(&broadcast, false)
			break
		}
	}
}

func (pool *Pool) SendHistoryToClient(client *Client) error {
	for _, message := range pool.History.GetMessages() {
		transport := NewMessageTransport(messages.ServerTextMessage, message)
		err := pool.BroadcastMessage(transport, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// LogInClient sends JWT to the client
func (pool *Pool) LogInClient(client *Client) error {
	jwt, err := client.User.GetJwt()
	if err != nil {
		return err
	}

	payload := messages.NewLoginResponseMessageObject(jwt)
	transport := NewMessageTransport(messages.LoginResponseMessage, payload)
	return pool.SendMessageToClient(transport, client)
}

func (pool *Pool) SendUsersQtyToClient(client *Client) error {
	qty := len(pool.getAuthorizedClients())
	payload := messages.NewParticipantsQtyUpdateObject(qty)
	transport := NewMessageTransport(messages.ParticipantsQtyUpdate, payload)
	return pool.SendMessageToClient(transport, client)
}

func (pool *Pool) BroadcastUsersQty() error {
	qty := len(pool.getAuthorizedClients())
	payload := messages.NewParticipantsQtyUpdateObject(qty)
	transport := NewMessageTransport(messages.ParticipantsQtyUpdate, payload)
	return pool.BroadcastMessage(transport, true)
}

func (pool *Pool) BroadcastMessage(message *MessageTransport, sendToAllClients bool) error {
	var clients []*Client

	if sendToAllClients {
		clients = pool.getAllClients()
	} else {
		clients = pool.getAuthorizedClients()
	}

	fmt.Println("Sending message to all clients in Pool")
	for _, client := range clients {
		if err := client.Conn.WriteJSON(message); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (pool *Pool) SendMessageToClient(message *MessageTransport, client *Client) error {
	fmt.Println("Sending message to client")
	if err := client.Conn.WriteJSON(message); err != nil {
		return err
	}
	return nil
}

func (pool *Pool) getAllClients() []*Client {
	var users []*Client
	for client, _ := range pool.Clients {
		users = append(users, client)
	}
	return users
}

func (pool *Pool) getAuthorizedClients() []*Client {
	var users []*Client
	for client, _ := range pool.Clients {
		if client.IsAuthorized() {
			users = append(users, client)
		}
	}
	return users
}

func (pool *Pool) LogOutClient(client *Client) error {
	client.User = nil
	pool.BroadcastServerMessage(fmt.Sprintf("%s Disconnected", client.User.Name))

	payload := messages.NewLogoutResponseMessageObject()
	transport := NewMessageTransport(messages.LogoutResponseMessage, payload)
	return pool.SendMessageToClient(transport, client)
}

func (pool *Pool) AddClientToPool(client *Client) {
	pool.Clients[client] = true
	fmt.Println("Size of Connection Pool: ", len(pool.Clients))
}

func (pool *Pool) DeleteClientFromPool(client *Client) error {
	delete(pool.Clients, client)
	fmt.Println("Size of Connection Pool: ", len(pool.Clients))
	pool.BroadcastUsersQty()
	if client.IsAuthorized() {
		return pool.BroadcastServerMessage(fmt.Sprintf("%s Disconnected", client.User.Name))
	}
	return nil
}

func (pool *Pool) BroadcastServerMessage(message string) error {
	payload := messages.NewServerTextMessage("Server", message)
	transport := NewMessageTransport(messages.ServerTextMessage, payload)
	pool.History.AddMessage(payload)
	return pool.BroadcastMessage(transport, false)
}
