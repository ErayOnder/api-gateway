package models

import "time"

// Conversation represents a conversation entity from chat-core
type Conversation struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	ModelName string    `json:"modelName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Messages  []Message `json:"messages,omitempty"`
}

// Message represents a message entity from chat-core
type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversationId"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	ResponseTimeMs *int      `json:"responseTimeMs,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

// IncomingUserMessage represents a message received from the UI via WebSocket
type IncomingUserMessage struct {
	ConversationID string `json:"conversationId"`
	UserMessage    string `json:"userMessage"`
}

// OutgoingBotMessage represents a response sent to the UI via WebSocket
type OutgoingBotMessage struct {
	BotMessage        *Message `json:"botMessage"`
	ConversationTitle string   `json:"conversationTitle,omitempty"` // Updated title if changed
}

// IncomingMessageRequest represents a request to send a message to chat-core
type IncomingMessageRequest struct {
	Content string `json:"content"`
}
