package main

import (
	"log"
	"os"

	"github.com/damiisdandy/contextual-chatbot/contexthandler"
	"github.com/damiisdandy/contextual-chatbot/fileprocessor"
	"github.com/damiisdandy/contextual-chatbot/screenshotprocessor"
	"github.com/joho/godotenv"
	"github.com/liushuangls/go-anthropic/v2"
)

func main() {
	// get environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Global context store

	// Parse conversation
	fileProcessor := fileprocessor.NewFileProcessor("damilola")
	messages, err := fileProcessor.Readfile("_chat.txt")
	if err != nil {
		log.Fatal(err)
	}

	// pass in the name of the person we are chatting with through chat logs
	contextStore := contexthandler.NewContextStore(fileProcessor.Peer)
	contextStore.AddMessages(messages)

	// get information from screenshot
	anthropicClient := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))
	screenshotProcessor := screenshotprocessor.NewScreenshotProcessor(anthropicClient)
	response, err := screenshotProcessor.ProcessImage("test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	messages, err = screenshotProcessor.ParseJSONString(response)
	if err != nil {
		log.Fatal(err)
	}
	contextStore.AddMessages(messages)

	contextStore.GeneratePromp("What do you think about our relationship?")

}
