package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) //upgrades http to websocket
	if err != nil {
		log.Println("websocket upgrade error:", err)
		return
	}
	defer conn.Close()
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
}
