package handlers

import (
	"api-gateway/internal/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// ConversationHandler handles HTTP conversation operations
type ConversationHandler struct {
	chatCoreClient *services.ChatCoreClient
}

// NewConversationHandler creates a new conversation handler
func NewConversationHandler(chatCoreClient *services.ChatCoreClient) *ConversationHandler {
	return &ConversationHandler{
		chatCoreClient: chatCoreClient,
	}
}

// CreateConversation creates a new conversation
func (h *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	conversation, err := h.chatCoreClient.CreateConversation()
	if err != nil {
		log.Printf("Error creating conversation: %v", err)
		http.Error(w, "Failed to create conversation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(conversation)
}

// ListConversations returns all conversations
func (h *ConversationHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	conversations, err := h.chatCoreClient.GetConversations()
	if err != nil {
		log.Printf("Error fetching conversations: %v", err)
		http.Error(w, "Failed to fetch conversations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

// GetConversation returns a single conversation with messages
func (h *ConversationHandler) GetConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["id"]

	conversation, err := h.chatCoreClient.GetConversation(conversationID)
	if err != nil {
		log.Printf("Error fetching conversation: %v", err)
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

// DeleteConversation deletes a conversation
func (h *ConversationHandler) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["id"]

	err := h.chatCoreClient.DeleteConversation(conversationID)
	if err != nil {
		log.Printf("Error deleting conversation: %v", err)
		http.Error(w, "Failed to delete conversation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
