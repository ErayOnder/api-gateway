package handlers

import (
	"api-gateway/internal/services"
	"api-gateway/pkg/models"
	"encoding/json"
	"log"
	"net/http"

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
		var wsMsg models.WebSocketMessage
		err = json.Unmarshal(message, &wsMsg)
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			h.sendErrorMessage(conn, "Invalid message format")
			continue
		}

		// Validate conversationId is provided
		if wsMsg.ConversationID == "" {
			log.Printf("Missing conversationId in message")
			h.sendErrorMessage(conn, "conversationId is required")
			continue
		}

		// Validate content
		if wsMsg.Content == "" {
			log.Printf("Missing content in message")
			h.sendErrorMessage(conn, "content is required")
			continue
		}

		// Call chat-core service to create message and get LLM response
		messagePair, err := h.chatCoreClient.SendMessage(wsMsg.ConversationID, wsMsg.Content)
		if err != nil {
			log.Printf("Error calling chat-core service: %v", err)
			h.sendErrorMessage(conn, "Failed to process your message")
			continue
		}

		// Send response back to client
		response := models.WebSocketResponse{
			Type:             "message",
			ConversationID:   messagePair.ConversationID,
			UserMessage:      messagePair.UserMessage,
			AssistantMessage: messagePair.AssistantMessage,
		}

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

	log.Println("Client disconnected from WebSocket")
}

// sendErrorMessage sends an error message to the client
func (h *WebSocketHandler) sendErrorMessage(conn *websocket.Conn, message string) {
	errorResponse := models.WebSocketResponse{
		Type:  "error",
		Error: message,
	}
	data, _ := json.Marshal(errorResponse)
	conn.WriteMessage(websocket.TextMessage, data)
}
