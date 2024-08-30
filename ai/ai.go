package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/liushuangls/go-anthropic/v2"
)

const MaxTokens = 1500
const AIModel = anthropic.ModelClaude3Haiku20240307
const ChatbotName = "Christie"

type AnthropicAI struct {
	Client *anthropic.Client
}

func NewAnthropicAI(apiKey string) *AnthropicAI {
	return &AnthropicAI{
		Client: anthropic.NewClient(apiKey),
	}
}

func (ai *AnthropicAI) executePrompt(ctx context.Context, messages []anthropic.Message) (string, error) {
	resp, err := ai.Client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:     AIModel,
		Messages:  messages,
		MaxTokens: MaxTokens,
	})
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			return "", fmt.Errorf("Messages error, type: %s, message: %s\n", e.Type, e.Message)
		} else {
			return "", fmt.Errorf("Messages error: %v\n", err)
		}
	}
	return resp.Content[0].GetText(), nil
}

func (ai *AnthropicAI) ExecutePrompt(ctx context.Context, prompt string) (string, error) {
	return ai.executePrompt(ctx, []anthropic.Message{
		anthropic.NewUserTextMessage(prompt),
	})
}

func (ai *AnthropicAI) AnalyseScreenshot(ctx context.Context, image []byte) (string, error) {
	var analyseScreenshotPrompt = fmt.Sprintf(
		`
	Provided is a whatsapp chat screenshot, I want you to only provide a json representatoin of this chat, it should contain a list of messages and each message should contain the following:
		1. sender
		2. content
		3. timestamp
	Also note the following:
	- Timestamp should be in the format of %s assume we are in the year 2024
	- The messages with the greenish chat bubble should be called the Reciever while the latter should be called the Sender
	- ONLY return the json stringified
	`, time.Now().Format("2006-01-02T15:04:05Z07:00"))

	return ai.executePrompt(ctx, []anthropic.Message{
		{
			Role: anthropic.RoleUser,
			Content: []anthropic.MessageContent{
				anthropic.NewImageMessageContent(anthropic.MessageContentImageSource{
					Type:      "base64",
					MediaType: "image/jpeg",
					Data:      image,
				}),
				anthropic.NewTextMessageContent(analyseScreenshotPrompt),
			},
		},
	})
}
