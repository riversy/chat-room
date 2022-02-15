package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/riversy/chat-room/server/pkg/websocket"
)

// serveWs
func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	buildPath := os.Getenv("STATIC_FILES_PATH")
	if buildPath == "" {
		buildPath = "../client/build"
	}
	fs := http.FileServer(http.Dir(buildPath))
	http.Handle("/", fs)
}

func main() {
	setupRoutes()
	fmt.Println("Listen http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
