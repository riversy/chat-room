package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/riversy/chat-room/server/pkg/websocket"
)

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		pool.ServeWs(w, r)
	})

	buildPath := os.Getenv("STATIC_FILES_PATH")
	if buildPath == "" {
		buildPath = "../client/build"
	}
	fs := http.FileServer(http.Dir(buildPath))
	http.Handle("/", fs)
}

func main() {
	godotenv.Load("../.env")
	setupRoutes()
	fmt.Println("Listen http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
