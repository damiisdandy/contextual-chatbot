package main

import (
	"fmt"

	"github.com/damiisdandy/contextual-chatbot/fileprocessor"
)

func main() {
	processor := fileprocessor.NewFileProcessor("damilola", "./conversations")
	conversations := processor.ParseConversations()

	for _, conversation := range conversations {
		fmt.Printf("Conversation: %s\n", conversation.Sender)
		fmt.Print(conversation.Messages[0].Content)
	}
}
