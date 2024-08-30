package tui

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/damiisdandy/contextual-chatbot/ai"
	ch "github.com/damiisdandy/contextual-chatbot/contexthandler"
	sp "github.com/damiisdandy/contextual-chatbot/screenshotprocessor"
)

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Magenta = "\033[35m"
const Cyan = "\033[36m"
const Gray = "\033[37m"
const White = "\033[97m"

type TUI struct {
	User                 string
	File                 string
	ProcessedScreenshots map[string]string
	ContextStore         *ch.ContextStore
	AIAgent              *ai.AnthropicAI
	ScreenshotProcessor  *sp.ScreenshotProcessor
}

func ProcessUserInput(description string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("> %s: ", description)
	fmt.Print(Green)
	input, err := reader.ReadString('\n')
	fmt.Print(Reset)
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(input)
}

func RenderErrorMessage(description string, err error) string {
	return fmt.Sprintf("%s\n %s. Please try again.\nError message:%s\n\n%s", Yellow, description, err, Reset)
}

func NewTUI() *TUI {
	fmt.Printf("> Hi, I'm %s your dating assistant. Let's start by uploading your WhatsApp conversation.\nPlease type in your name (preferably your name on Whatsapp) and the name of the conversation file (The conversation file should be in the 'conversations' folder)\n\n", ai.ChatbotName)

	user := ProcessUserInput("Your name")
	file := ProcessUserInput("Conversation file name")

	fmt.Printf("\n> Thank you for providing the conversation file. I'm currently processing the file. Please wait...")
	return &TUI{
		User:                 user,
		File:                 file,
		ProcessedScreenshots: map[string]string{},
	}
}

func (t *TUI) AddScreenshot(filename, content string) {
	t.ProcessedScreenshots[filename] = content
}

func (t *TUI) ScreenshotExists(filename string) bool {
	_, ok := t.ProcessedScreenshots[filename]
	return ok
}

func (t *TUI) HandleAskQuestion(additionalInfo string) {
	question := ProcessUserInput("Type your question")
	fmt.Printf("\nThinking...\n")
	prompt := t.ContextStore.GeneratePrompt(question)
	response, err := t.AIAgent.ExecutePrompt(context.Background(), prompt)
	if err != nil {
		RenderErrorMessage("Problem connecting to the AI", err)
		return
	}
	fmt.Printf("\n> %s%s%s\n\n", Blue, response, Reset)
	t.ContextStore.AddQuestion(question)
	return
}

func (t *TUI) HandleAskQuestionWithScreenshot() {
	filename := ProcessUserInput("Please input the screenshot file name (The screenshot should be in the 'screenshots' folder)")
	if !t.ScreenshotExists(filename) {
		fmt.Printf("\nProcessing screenshot...\n")
		// Process screenshot
		jsonString, err := t.ScreenshotProcessor.ProcessImage(filename)
		if err != nil {
			RenderErrorMessage("Problem processing the screenshot", err)
			return
		}
		messages, err := t.ScreenshotProcessor.ParseJSONString(jsonString)
		if err != nil {
			RenderErrorMessage("Problem parsing the screenshot", err)
			return
		}
		t.ContextStore.AddMessages(messages, ch.MessageSourceScreenshot)
		t.AddScreenshot(filename, jsonString)
	} else {
		fmt.Printf("%s\nLooks like I've already looked into screenshot %q. You can ask your question%s\n\n", Yellow, filename, Reset)
	}
	// Ask question
	t.HandleAskQuestion("\n based on the most recent screenshot sent, reply based on the recent call logs of source " + ch.MessageSourceScreenshot)
}
