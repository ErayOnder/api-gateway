package services

import (
	"api-gateway/pkg/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ChatCoreClient handles communication with the chat-core service
type ChatCoreClient struct {
	baseURL string
	client  *http.Client
}

// NewChatCoreClient creates a new chat-core service client
func NewChatCoreClient(baseURL string) *ChatCoreClient {
	return &ChatCoreClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// CreateConversation creates a new conversation
func (c *ChatCoreClient) CreateConversation() (*models.Conversation, error) {
	resp, err := c.client.Post(
		c.baseURL+"/conversations",
		"application/json",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call chat-core service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("chat-core service returned status code: %d", resp.StatusCode)
	}

	var conversation models.Conversation
	err = json.NewDecoder(resp.Body).Decode(&conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &conversation, nil
}

// GetConversations retrieves all conversations
func (c *ChatCoreClient) GetConversations() ([]models.Conversation, error) {
	resp, err := c.client.Get(c.baseURL + "/conversations")
	if err != nil {
		return nil, fmt.Errorf("failed to call chat-core service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chat-core service returned status code: %d", resp.StatusCode)
	}

	var conversations []models.Conversation
	err = json.NewDecoder(resp.Body).Decode(&conversations)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return conversations, nil
}

// GetConversation retrieves a single conversation with messages
func (c *ChatCoreClient) GetConversation(conversationID string) (*models.Conversation, error) {
	url := fmt.Sprintf("%s/conversations/%s", c.baseURL, conversationID)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call chat-core service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chat-core service returned status code: %d", resp.StatusCode)
	}

	var conversation models.Conversation
	err = json.NewDecoder(resp.Body).Decode(&conversation)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &conversation, nil
}

// DeleteConversation deletes a conversation
func (c *ChatCoreClient) DeleteConversation(conversationID string) error {
	url := fmt.Sprintf("%s/conversations/%s", c.baseURL, conversationID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call chat-core service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("chat-core service returned status code: %d", resp.StatusCode)
	}

	return nil
}

// SendMessage sends a user message and gets an AI response from chat-core service
func (c *ChatCoreClient) SendMessage(conversationID, content string) (*models.Message, error) {
	requestBody := models.IncomingMessageRequest{
		Content: content,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/conversations/%s/messages/chat", c.baseURL, conversationID)
	resp, err := c.client.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call chat-core service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("chat-core service returned status code: %d", resp.StatusCode)
	}

	var response models.Message
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
