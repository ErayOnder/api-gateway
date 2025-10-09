package handlers

import (
	"api-gateway/internal/services"
	"api-gateway/pkg/models"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	chatCoreClient *services.ChatCoreClient
	writeMutex     sync.Mutex // Protects concurrent writes to WebSocket
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(chatCoreClient *services.ChatCoreClient) *WebSocketHandler {
	return &WebSocketHandler{
		chatCoreClient: chatCoreClient,
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

	log.Println("Client connected to WebSocket")

	for {
		// Read message from client
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		log.Printf("Received message: %s", message)

		// Parse WebSocket message
		var incomingUserMessage models.IncomingUserMessage
		err = json.Unmarshal(message, &incomingUserMessage)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			h.sendErrorMessage(conn, "Invalid message format")
			continue
		}

		// Validate conversationId is provided
		if incomingUserMessage.ConversationID == "" {
			log.Printf("Missing conversationId in message")
			h.sendErrorMessage(conn, "conversationId is required")
			continue
		}

		// Validate content
		if incomingUserMessage.UserMessage == "" {
			log.Printf("Missing content in message")
			h.sendErrorMessage(conn, "content is required")
			continue
		}

		// Handle message asynchronously in a goroutine
		// This allows processing multiple messages concurrently
		go h.processMessage(conn, incomingUserMessage)
	}

	log.Println("Client disconnected from WebSocket")
}

// processMessage handles a single message asynchronously
func (h *WebSocketHandler) processMessage(conn *websocket.Conn, incomingUserMessage models.IncomingUserMessage) {
	// Call chat-core service to create message and get LLM response
	botMessage, err := h.chatCoreClient.SendMessage(incomingUserMessage.ConversationID, incomingUserMessage.UserMessage)
	if err != nil {
		log.Printf("Error calling chat-core service: %v", err)
		h.sendErrorMessage(conn, "Failed to process your message")
		return
	}

	// Send response back to client
	response := models.OutgoingBotMessage{
		BotMessage: botMessage,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	// Use mutex to safely write to WebSocket from multiple goroutines
	h.writeMutex.Lock()
	err = conn.WriteMessage(websocket.TextMessage, responseData)
	h.writeMutex.Unlock()

	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

// sendErrorMessage sends an error message to the WebSocket client
func (h *WebSocketHandler) sendErrorMessage(conn *websocket.Conn, errorMsg string) {
	errorResponse := map[string]string{
		"error": errorMsg,
	}

	data, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Error marshaling error response: %v", err)
		return
	}

	// Use mutex to safely write to WebSocket from multiple goroutines
	h.writeMutex.Lock()
	err = conn.WriteMessage(websocket.TextMessage, data)
	h.writeMutex.Unlock()

	if err != nil {
		log.Printf("Error sending error message: %v", err)
	}
}
