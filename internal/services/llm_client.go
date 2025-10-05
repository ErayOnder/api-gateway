package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"api-gateway/pkg/models"
)

// LLMClient handles communication with the LLM service
type LLMClient struct {
	baseURL string
	client  *http.Client
}

// NewLLMClient creates a new LLM service client
func NewLLMClient(baseURL string) *LLMClient {
	return &LLMClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Generate sends a message to the LLM service and returns the response
func (c *LLMClient) Generate(message string) (*models.LLMResponse, error) {
	requestBody := models.MessageRequest{
		Message: message,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(
		c.baseURL+"/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("LLM service returned status code: %d", resp.StatusCode)
	}

	var llmResponse models.LLMResponse
	err = json.NewDecoder(resp.Body).Decode(&llmResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	return &llmResponse, nil
}
