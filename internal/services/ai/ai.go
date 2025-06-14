package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var AvailableModels = []string{
	"gpt-4o",
	"gpt-4.1",
	"o1",
	"o4-mini",
	"deepseek-r1-0528",
	"deepseek-v3-0324",
	"grok-3",
}

type AzureAiRequest struct {
	Stream   bool      `json:"stream" default:"false"`
	Model    string    `json:"model" binding:"required"`
	Messages []Message `json:"messages" binding:"required"`
}

type Message struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type AzureAiResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func Request(request AzureAiRequest) (AzureAiResponse, error) {
	endpoint := os.Getenv("AZURE_AI_ENDPOINT") + "/models/chat/completions"
	key := os.Getenv("AZURE_AI_KEY")

	jsonData, err := json.Marshal(request)
	if err != nil {
		return AzureAiResponse{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return AzureAiResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", key)

	resp, err := client.Do(req)
	if err != nil {
		return AzureAiResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AzureAiResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return AzureAiResponse{}, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response AzureAiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return AzureAiResponse{}, err
	}

	return response, nil
}

func RequestStream(request AzureAiRequest) (*http.Response, error) {
	endpoint := os.Getenv("AZURE_AI_ENDPOINT") + "/models/chat/completions"
	key := os.Getenv("AZURE_AI_KEY")

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", key)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
