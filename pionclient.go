// peer_connection.go

package main

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func poinClient() {
	// Get signaling server address from environment variable
	serverAddr := "localhost:30001" //os.Getenv("SIGNALING_SERVER_ADDRESS")
	if serverAddr == "" {
		log.Fatal("SIGNALING_SERVER_ADDRESS environment variable is not set")
	}

	// Connect to the signaling server
	log.Printf("Connecting to signaling server at %s\n", serverAddr)
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+serverAddr+"/ws", nil)
	if err != nil {
		log.Fatal("Error connecting to signaling server:", err)
	}
	defer conn.Close()

	log.Println("Connected to signaling server")
	// Example: Sending a message to the signaling server
	if err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, signaling server!")); err != nil {
		log.Println("Error sending message to signaling server:", err)
	}

	// Create a new WebRTC peer connection configuration
	config := webrtc.Configuration{}

	// Create a new WebRTC peer connection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal("Error creating peer connection:", err)
	}
	defer peerConnection.Close()

	log.Println("WebRTC peer connection created successfully")

	// Handle signaling messages from the signaling server
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from signaling server:", err)
				break
			}

			// Handle the received signaling message
			log.Printf("Received signaling message from server: %s\n", message)

			// Implement logic to parse and handle the signaling message here
		}
	}()

	// Wait indefinitely to keep the program running
	select {}
}
