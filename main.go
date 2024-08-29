package main

import (
	"fmt"
	"log"
	"os"

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

	anthropicClient := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))

	screenshotProcessor := screenshotprocessor.NewScreenshotProcessor(anthropicClient)
	response, err := screenshotProcessor.ProcessImage("test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

}
