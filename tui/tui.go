package tui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/damiisdandy/contextual-chatbot/ai"
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
