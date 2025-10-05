package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"api-gateway/internal/services"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	llmClient *services.LLMClient
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(llmClient *services.LLMClient) *WebSocketHandler {
	return &WebSocketHandler{
		llmClient: llmClient,
	}
}

// Handle processes WebSocket connections
func (h *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		log.Printf("Received message: %s", message)

		// Call LLM service
		response, err := h.llmClient.Generate(string(message))
		if err != nil {
			log.Printf("Error calling LLM service: %v", err)
			h.sendErrorMessage(conn, "Failed to process your message")
			continue
		}

		// Send response back to client
		responseData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, responseData)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}

	log.Println("Client disconnected")
}

// sendErrorMessage sends an error message to the client
func (h *WebSocketHandler) sendErrorMessage(conn *websocket.Conn, message string) {
	errorResponse := map[string]string{
		"error": message,
	}
	data, _ := json.Marshal(errorResponse)
	conn.WriteMessage(websocket.TextMessage, data)
}
