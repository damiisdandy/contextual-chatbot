package screenshotprocessor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/damiisdandy/contextual-chatbot/contexthandler"
	"github.com/liushuangls/go-anthropic/v2"
)

var ProcessScreenshotPrompt = fmt.Sprintf(`Provided is a whatsapp chat screenshot, I want you to only provide a json representatoin of this chat,
																 it should contain a list of messages and each message should contain the following:
																	1. sender
																	2. content
																	3. timestamp
																Also note the following:
																- Timestamp should be in the format of %s assume we are in the year 2024
															  - The messages with the greenish chat bubble should be called the Reciever while the latter should be called the Sender
																- ONLY return the json stringified
																`, time.Now().Format("2006-01-02T15:04:05Z07:00"))

type ScreenshotProcessor struct {
	ScreenshopPath  string
	AnthropicClient *anthropic.Client
}

func NewScreenshotProcessor(anthropicClient *anthropic.Client) *ScreenshotProcessor {
	return &ScreenshotProcessor{
		ScreenshopPath:  filepath.Join("./screenshots"),
		AnthropicClient: anthropicClient,
	}
}

func (sp *ScreenshotProcessor) ProcessImage(imagePath string) (string, error) {
	anthropicClient := sp.AnthropicClient

	imageFile, err := os.Open(filepath.Join(sp.ScreenshopPath, imagePath))
	if err != nil {
		return "", errors.New("Error opening image file")
	}
	imageData, err := io.ReadAll(imageFile)
	if err != nil {
		return "", errors.New("Error reading image file")
	}

	resp, err := anthropicClient.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Haiku20240307,
		Messages: []anthropic.Message{
			{
				Role: anthropic.RoleUser,
				Content: []anthropic.MessageContent{
					anthropic.NewImageMessageContent(anthropic.MessageContentImageSource{
						Type:      "base64",
						MediaType: "image/jpeg",
						Data:      imageData,
					}),
					anthropic.NewTextMessageContent(ProcessScreenshotPrompt),
				},
			},
		},
		MaxTokens: 1500,
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			return "", fmt.Errorf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			return "", fmt.Errorf("Messages error: %v\n", err)
		}
	}
	return resp.Content[0].GetText(), nil
}

func (s *ScreenshotProcessor) ParseJSONString(jsonString string) ([]contexthandler.Message, error) {
	messages := []contexthandler.Message{}

	// ensure we parse only the JSON stringified
	start := strings.Index(jsonString, "[")
	end := strings.LastIndex(jsonString, "]")

	err := json.Unmarshal([]byte(jsonString[start:end+1]), &messages)
	if err != nil {
		return nil, fmt.Errorf("Error parsing AI response: %s", err)
	}

	return messages, nil
}
