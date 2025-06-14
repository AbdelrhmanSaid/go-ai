package handlers

import (
	"bufio"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/AbdelrhmanSaid/go-ai/internal/services/ai"
	"github.com/gin-gonic/gin"
)

type ChatCompletionsRequest struct {
	Model    string    `json:"model" binding:"required"`
	Messages []Message `json:"messages" binding:"required"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (r *ChatCompletionsRequest) Validate() error {
	if !slices.Contains(ai.AvailableModels, r.Model) {
		return errors.New("invalid model")
	}

	if len(r.Messages) == 0 {
		return errors.New("messages cannot be empty")
	}

	for _, message := range r.Messages {
		if message.Role != "user" && message.Role != "assistant" {
			return errors.New("invalid role")
		}
	}

	return nil
}

func ChatCompletions(c *gin.Context) {
	var request ChatCompletionsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Map handlers.Message to ai.Message
	aiMessages := make([]ai.Message, len(request.Messages))
	for i, m := range request.Messages {
		aiMessages[i] = ai.Message{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	if request.Stream {
		// Handle streaming response
		resp, err := ai.RequestStream(ai.AzureAiRequest{
			Model:    request.Model,
			Messages: aiMessages,
			Stream:   true,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		defer resp.Body.Close()

		// Set headers for SSE
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")

		// Stream the response
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "streaming not supported",
			})
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			// Remove "data: " prefix if present
			line = strings.TrimPrefix(line, "data: ")

			// Check for end of stream
			if line == "[DONE]" {
				c.SSEvent("", "[DONE]")
				flusher.Flush()
				return
			}

			// Forward the chunk to client
			c.SSEvent("", line)
			flusher.Flush()
		}
	} else {
		// Handle non-streaming response (existing logic)
		response, err := ai.Request(ai.AzureAiRequest{
			Model:    request.Model,
			Messages: aiMessages,
			Stream:   false,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": response.Choices[0].Message.Content,
		})
	}
}
