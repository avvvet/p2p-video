// peer_connection.go

package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var peerConnection *webrtc.PeerConnection

func pionClient() {
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
			// Implement logic to parse and handle the signaling message here
			handleSignalingMessage(message)
		}
	}()

	// Wait indefinitely to keep the program running
	select {}
}

// Function to handle signaling messages
// Function to handle signaling messages
func handleSignalingMessage(message []byte) {
	// Parse the message and determine its type
	// For now, we'll just log the received message
	log.Printf("Received signaling message: %s\n", message)

	// Parse the message as a JSON object
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Println("Error parsing signaling message:", err)
		return
	}

	// Check the message type
	messageType, ok := msg["type"].(string)
	if !ok {
		log.Println("Invalid message type")
		return
	}

	// Handle SDP offers/answers and ICE candidates based on the message type
	switch messageType {
	case "offer":
		// Handle SDP offer
		handleSDPOffer(msg)
	case "answer":
		// Handle SDP answer
		handleSDPAnswer(msg)
	case "candidate":
		// Handle ICE candidate
		handleICECandidate(msg)
	default:
		log.Println("Unsupported message type:", messageType)
	}
}

// Function to handle ICE candidate
// Function to handle ICE candidate
func handleICECandidate(msg map[string]interface{}) {
	// Extract the ICE candidate information from the message
	candidate := webrtc.ICECandidateInit{
		Candidate: msg["candidate"].(string),
	}

	sdpMid, sdpMidExists := msg["sdpMid"].(string)
	if sdpMidExists {
		candidate.SDPMid = &sdpMid
	}

	sdpMLineIndex, sdpMLineIndexExists := msg["sdpMLineIndex"].(float64)
	if sdpMLineIndexExists {
		sdpMLineIndexInt := uint16(sdpMLineIndex)
		candidate.SDPMLineIndex = &sdpMLineIndexInt
	}

	// Add the ICE candidate to the peer connection
	if err := peerConnection.AddICECandidate(candidate); err != nil {
		log.Println("Error adding ICE candidate:", err)
		return
	}

	log.Println("ICE candidate added successfully")
}

// Function to handle SDP offer
func handleSDPOffer(msg map[string]interface{}) {
	// Extract the SDP offer from the message
	sdpOffer, ok := msg["sdp"].(string)
	if !ok {
		log.Println("Invalid SDP offer format")
		return
	}

	// Create a new WebRTC session description
	offer := webrtc.SessionDescription{
		SDP:  sdpOffer,
		Type: webrtc.SDPTypeOffer,
	}

	// Set the remote description of the peer connection
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		log.Println("Error setting remote description:", err)
		return
	}

	log.Println("SDP offer set successfully")
}

// Function to handle SDP answer
func handleSDPAnswer(msg map[string]interface{}) {
	// Extract the SDP answer from the message
	sdpAnswer, ok := msg["sdp"].(string)
	if !ok {
		log.Println("Invalid SDP answer format")
		return
	}

	// Create a new WebRTC session description
	answer := webrtc.SessionDescription{
		SDP:  sdpAnswer,
		Type: webrtc.SDPTypeAnswer,
	}

	// Set the remote description of the peer connection
	if err := peerConnection.SetRemoteDescription(answer); err != nil {
		log.Println("Error setting remote description:", err)
		return
	}

	log.Println("SDP answer set successfully")
}
