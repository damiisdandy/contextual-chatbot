package main

import (
	"fmt"

	"github.com/damiisdandy/contextual-chatbot/fileprocessor"
)

func main() {
	processor := fileprocessor.NewFileProcessor("damilola", "_chat.txt")
	conversation, err := processor.ParseConversations()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Conversation: %s\n", conversation.Sender)
	fmt.Print(conversation.Messages[0].Content)
}
