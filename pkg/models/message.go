package models

// MessageRequest represents the request payload sent to the LLM service
type MessageRequest struct {
	Message string `json:"message"`
}

// LLMResponse represents the response from the LLM service
type LLMResponse struct {
	Response string `json:"response"`
}
