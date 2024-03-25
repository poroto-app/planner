package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ChatCompletionClient struct {
	apiKey string
}

func NewChatCompletionClient() (*ChatCompletionClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}

	return &ChatCompletionClient{
		apiKey: apiKey,
	}, nil
}

/*
SEE: https://platform.openai.com/docs/api-reference/chat
*/

const (
	ModelGPT3Turbo = "gpt-3.5-turbo"
)

type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
	N        *int                    `json:"n,omitempty"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Choices []ChatCompletionChoice       `json:"choices"`
	Error   *ChatCompletionErrorResponse `json:"error,omitempty"`
}

type ChatCompletionErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type ChatCompletionChoice struct {
	Index   int                   `json:"index"`
	Message ChatCompletionMessage `json:"message"`
}

func (c *ChatCompletionClient) Complete(request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	body := request

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, fmt.Errorf("error while encoding body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", &buf)
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while sending request: %v", err)
	}
	defer resp.Body.Close()

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error while decoding response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		if response.Error == nil {
			return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
		}

		return nil, fmt.Errorf("unexpected status code(%v): %v", resp.StatusCode, response.Error)
	}

	return &response, nil
}
