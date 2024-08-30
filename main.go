package main

import (
	"fmt"
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

	// Parse conversation
	fileProcessor := fileprocessor.NewFileProcessor("damilola")
	messages, err := fileProcessor.Readfile("_chat.txt")
	if err != nil {
		log.Fatal(err)
	}

	// pass in the name of the person we are chatting with through chat logs
	contextStore := contexthandler.NewContextStore(fileProcessor.Peer)
	contextStore.AddMessages(messages, contexthandler.MessageSourceLogs)

	// get information from screenshot
	anthropicClient := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))
	screenshotProcessor := screenshotprocessor.NewScreenshotProcessor(anthropicClient)
	response, err := screenshotProcessor.ProcessImage("screenshot.jpg")
	fmt.Print(response)
	if err != nil {
		log.Fatal(err)
	}
	messages, err = screenshotProcessor.ParseJSONString(response)
	if err != nil {
		log.Fatal(err)
	}
	contextStore.AddMessages(messages, contexthandler.MessageSourceScreenshot)

	// promote := contextStore.GeneratePromp("What do you think about our relationship? based on the most recent screenshot sent")
	// fmt.Print(promote)

	// resp, err := anthropicClient.CreateMessages(context.Background(), anthropic.MessagesRequest{
	// 	Model: anthropic.ModelClaude3Haiku20240307,
	// 	Messages: []anthropic.Message{
	// 		anthropic.NewUserTextMessage(promote),
	// 	},
	// 	MaxTokens: 1500,
	// })
	// if err != nil {
	// 	var e *anthropic.APIError
	// 	if errors.As(err, &e) {
	// 		fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
	// 	} else {
	// 		fmt.Printf("Messages error: %v\n", err)
	// 	}
	// 	return
	// }
	// fmt.Println(resp.Content[0].GetText())

}
