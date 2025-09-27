package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type LLMResponse struct {
	Response string `json:"response"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		log.Printf("Received message: %s", message)

		response, err := callLLMService(message)
		if err != nil {
			log.Printf("Error calling LLM service: %v", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func callLLMService(message []byte) ([]byte, error) {
	requestBody := map[string]any{
		"message": string(message),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	llmServiceURL := os.Getenv("LLM_SERVICE_URL")
	if llmServiceURL == "" {
		llmServiceURL = "http://localhost:3000"
	}

	resp, err := http.Post(llmServiceURL+"/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var llmResponse LLMResponse
	err = json.NewDecoder(resp.Body).Decode(&llmResponse)
	if err != nil {
		return nil, err
	}

	responseData, err := json.Marshal(llmResponse)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("API Gateway starting on :8080")
	fmt.Println("WebSocket endpoint available at ws://localhost:8080/ws")

	log.Fatal(http.ListenAndServe(":8080", nil))
}