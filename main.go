package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/damiisdandy/contextual-chatbot/ai"
	"github.com/damiisdandy/contextual-chatbot/contexthandler"
	"github.com/damiisdandy/contextual-chatbot/fileprocessor"
	sp "github.com/damiisdandy/contextual-chatbot/screenshotprocessor"
	"github.com/joho/godotenv"
)

func main() {
	// get environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	aiAgent := ai.NewAnthropicAI(os.Getenv("ANTHROPIC_API_KEY"))
	screenshotProcessor := sp.NewScreenshotProcessor(aiAgent)

	// Parse conversation
	fileProcessor := fileprocessor.NewFileProcessor("damilola")
	messages, err := fileProcessor.Readfile("_chat.txt")
	if err != nil {
		log.Fatal(err)
	}

	// pass in the name of the person we are chatting with through chat logs
	contextStore := contexthandler.NewContextStore("damilola", fileProcessor.Peer)
	contextStore.AddMessages(messages, contexthandler.MessageSourceLogs)

	// get information from screenshot

	response, err := screenshotProcessor.ProcessImage("screenshot.jpg")
	if err != nil {
		log.Fatal(err)
	}
	messages, err = screenshotProcessor.ParseJSONString(response)
	if err != nil {
		log.Fatal(err)
	}
	contextStore.AddMessages(messages, contexthandler.MessageSourceScreenshot)

	prompt := contextStore.GeneratePromp("What do you think about our relationship? based on the most recent screenshot sent")

	response, err = aiAgent.ExecutePrompt(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(response)
}
