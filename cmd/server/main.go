package main

import (
	"fmt"
	"log"
	"net/http"

	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/services"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize chat-core client
	chatCoreClient := services.NewChatCoreClient(cfg.ChatCoreURL)

	// Initialize handlers
	wsHandler := handlers.NewWebSocketHandler(chatCoreClient)
	conversationHandler := handlers.NewConversationHandler(chatCoreClient)

	// Setup router
	router := mux.NewRouter()

	// Apply CORS and logging middleware globally
	router.Use(middleware.CORS)
	router.Use(middleware.Logging)

	// WebSocket endpoint
	router.HandleFunc("/ws", wsHandler.Handle)

	// HTTP API endpoints for conversations
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/conversations", conversationHandler.ListConversations).Methods("GET", "OPTIONS")
	api.HandleFunc("/conversations", conversationHandler.CreateConversation).Methods("POST", "OPTIONS")
	api.HandleFunc("/conversations/{id}", conversationHandler.GetConversation).Methods("GET", "OPTIONS")
	api.HandleFunc("/conversations/{id}", conversationHandler.DeleteConversation).Methods("DELETE", "OPTIONS")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	addr := ":" + cfg.ServerPort
	fmt.Printf("API Gateway starting on %s\n", addr)
	fmt.Printf("WebSocket endpoint: ws://localhost:%s/ws\n", cfg.ServerPort)
	fmt.Printf("HTTP API: http://localhost:%s/api\n", cfg.ServerPort)
	fmt.Printf("Chat-Core Service URL: %s\n", cfg.ChatCoreURL)

	log.Fatal(http.ListenAndServe(addr, router))
}
