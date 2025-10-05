package main

import (
	"fmt"
	"log"
	"net/http"

	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"api-gateway/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize LLM client
	llmClient := services.NewLLMClient(cfg.LLMServiceURL)

	// Initialize handlers
	wsHandler := handlers.NewWebSocketHandler(llmClient)

	// Setup routes with middleware
	http.HandleFunc("/ws", middleware.LoggingMiddleware(
		middleware.CORSMiddleware(wsHandler.Handle),
	))

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	addr := ":" + cfg.ServerPort
	fmt.Printf("API Gateway starting on %s\n", addr)
	fmt.Printf("WebSocket endpoint available at ws://localhost:%s/ws\n", cfg.ServerPort)
	fmt.Printf("LLM Service URL: %s\n", cfg.LLMServiceURL)

	log.Fatal(http.ListenAndServe(addr, nil))
}
