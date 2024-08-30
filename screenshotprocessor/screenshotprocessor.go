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

	"github.com/damiisdandy/contextual-chatbot/ai"
	"github.com/damiisdandy/contextual-chatbot/contexthandler"
)

type ScreenshotProcessor struct {
	ScreenshopPath string
	AIAgent        *ai.AnthropicAI
}

func NewScreenshotProcessor(aiAgent *ai.AnthropicAI) *ScreenshotProcessor {
	return &ScreenshotProcessor{
		ScreenshopPath: filepath.Join("./screenshots"),
		AIAgent:        aiAgent,
	}
}

func (sp *ScreenshotProcessor) ProcessImage(imagePath string) (string, error) {
	imageFile, err := os.Open(filepath.Join(sp.ScreenshopPath, imagePath))
	if err != nil {
		return "", errors.New("Error opening image file")
	}
	imageData, err := io.ReadAll(imageFile)
	if err != nil {
		return "", errors.New("Error reading image file")
	}
	return sp.AIAgent.AnalyseScreenshot(context.Background(), imageData)
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
