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

	/* lets continuesly handle websocket message */
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read error:", err)
			break
		}

		log.Println(string(message))
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Signaling server listening...")
	log.Fatal(http.ListenAndServe(":30001", nil))
}
