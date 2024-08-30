package main

import (
	"context"
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
		fmt.Printf("%s\nLooks like there was a problem reading the file. Please try again.\nError message: %s\n%s", tui.Yellow, err, tui.Reset)
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

	for {
		fmt.Printf("> Please choose from the following options:\n 1. Ask a question\n 2. Ask a question with an attached screenshot\n\n")
		option := tui.ProcessUserInput("Enter option")
		switch option {
		case "1":
			question := tui.ProcessUserInput("Type your question")
			fmt.Printf("\nThinking...\n")
			prompt := contextStore.GeneratePrompt(question)
			response, err := aiAgent.ExecutePrompt(context.Background(), prompt)
			if err != nil {
				fmt.Printf("%s\nProblem connecting to the AI. Please try again.\nError message: %s\n\n%s", tui.Yellow, err, tui.Reset)
				continue
			}
			fmt.Printf("\n> %s%s%s\n\n", tui.Blue, response, tui.Reset)
			contextStore.AddQuestion(question)
			continue
		case "2":
			filename := tui.ProcessUserInput("Please input the screenshot file name (The screenshot should be in the 'screenshots' folder)")
			if !root.ScreenshotExists(filename) {
				fmt.Printf("\nProcessing screenshot...\n")
				// Process screenshot
				jsonString, err := screenshotProcessor.ProcessImage(filename)
				if err != nil {
					fmt.Printf("%s\nProblem processing the screenshot. Please try again.\nError message: %s\n\n%s", tui.Yellow, err, tui.Reset)
					continue
				}
				messages, err := screenshotProcessor.ParseJSONString(jsonString)
				if err != nil {
					fmt.Printf("%s\nProblem parsing the screenshot. Please try again.\nError message: %s\n\n%s", tui.Yellow, err, tui.Reset)
					continue
				}
				contextStore.AddMessages(messages, ch.MessageSourceScreenshot)
				root.AddScreenshot(filename, jsonString)
			} else {
				fmt.Printf("%s\nLooks like I've already looked into screenshot %q. You can ask your question%s\n\n", tui.Yellow, filename, tui.Reset)
			}

			// Ask question
			question := tui.ProcessUserInput("Type your question")
			fmt.Printf("\nThinking...\n")
			prompt := contextStore.GeneratePrompt(question + "\n based on the most recent screenshot sent")
			response, err := aiAgent.ExecutePrompt(context.Background(), prompt)
			if err != nil {
				fmt.Printf("%s\nProblem connecting to the AI. Please try again.\nError message: %s\n\n%s", tui.Yellow, err, tui.Reset)
				continue
			}
			fmt.Printf("\n> %s%s%s\n\n", tui.Blue, response, tui.Reset)
			contextStore.AddQuestion(question)
			continue
		case "exit":
			fmt.Printf("\nThank you for using the contextual chatbot. Have a great day!\n\n")
			os.Exit(0)
		default:
			fmt.Printf("%s\nInvalid option %q, choose between 1 and 2. Please try again.\n\n%s", tui.Yellow, option, tui.Reset)
			continue
		}
	}
}
