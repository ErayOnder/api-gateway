package models

import "time"

// Conversation represents a conversation entity from chat-core
type Conversation struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	ModelName    string    `json:"modelName"`
	SystemPrompt *string   `json:"systemPrompt,omitempty"`
	IsArchived   bool      `json:"isArchived"`
	IsPinned     bool      `json:"isPinned"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Messages     []Message `json:"messages,omitempty"`
}

// Message represents a message entity from chat-core
type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversationId"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	TokenCount     *int      `json:"tokenCount,omitempty"`
	ResponseTimeMs *int      `json:"responseTimeMs,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

// MessagePairResponse contains both user and assistant messages
type MessagePairResponse struct {
	UserMessage      *Message `json:"userMessage"`
	AssistantMessage *Message `json:"assistantMessage"`
	ConversationID   string   `json:"conversationId"`
}

// WebSocketMessage represents a message received from the UI via WebSocket
type WebSocketMessage struct {
	ConversationID string `json:"conversationId"` // Required
	Content        string `json:"content"`        // Required
}

// WebSocketResponse represents a response sent to the UI via WebSocket
type WebSocketResponse struct {
	Type             string   `json:"type"` // "message", "error", "conversation_created"
	ConversationID   string   `json:"conversationId,omitempty"`
	UserMessage      *Message `json:"userMessage,omitempty"`
	AssistantMessage *Message `json:"assistantMessage,omitempty"`
	Error            string   `json:"error,omitempty"`
}

// CreateConversationRequest represents the request to create a conversation
type CreateConversationRequest struct {
	Title        string  `json:"title"`
	ModelName    string  `json:"modelName"`
	SystemPrompt *string `json:"systemPrompt,omitempty"`
}

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	Content string `json:"content"`
}
