package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/damiisdandy/contextual-chatbot/ai"
	ch "github.com/damiisdandy/contextual-chatbot/contexthandler"
	fp "github.com/damiisdandy/contextual-chatbot/fileprocessor"
	sp "github.com/damiisdandy/contextual-chatbot/screenshotprocessor"
	"github.com/damiisdandy/contextual-chatbot/tui"
	"github.com/joho/godotenv"
)

func main() {
	// get environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Terminal UI
	root := tui.NewTUI()

	// File processor
	fileProcessor := fp.NewFileProcessor(root.User)
	messages, err := fileProcessor.Readfile(filepath.Join("./conversations", root.File))
	if err != nil {
		tui.RenderErrorMessage("Looks like there was a problem reading the file", err)
		return
	}

	// Context store
	contextStore := ch.NewContextStore(root.User, fileProcessor.Peer) // Peer is the name of the person we are chatting with
	contextStore.AddMessages(messages, ch.MessageSourceLogs)

	fmt.Printf("\n> Your conversation file is processed successfully. Now you can ask questions.\n\n")

	// AI Agent
	aiAgent := ai.NewAnthropicAI(os.Getenv("ANTHROPIC_API_KEY"))

	// Screenshot processor
	screenshotProcessor := sp.NewScreenshotProcessor(aiAgent)

	// Assign dependencies
	root.AIAgent = aiAgent
	root.ScreenshotProcessor = screenshotProcessor
	root.ContextStore = contextStore

	for {
		fmt.Printf("> Please choose from the following options:\n 1. Ask a question\n 2. Ask a question with an attached screenshot\n\n")
		option := tui.ProcessUserInput("Enter option")
		switch option {
		case "1":
			root.HandleAskQuestion("")
			continue
		case "2":
			root.HandleAskQuestionWithScreenshot()
		case "exit":
			fmt.Printf("\nThank you for using the contextual chatbot. Have a great day!\n\n")
			os.Exit(0)
		default:
			fmt.Printf("%s\nInvalid option %q, choose between 1 and 2. Please try again.\n\n%s", tui.Yellow, option, tui.Reset)
			continue
		}
	}
}
