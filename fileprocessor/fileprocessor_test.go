package fileprocessor

import (
	"fmt"
	"testing"
	"time"
)

var fileProcessor FileProcessor = FileProcessor{
	Reciever: "damilola",
}

func TestParseMessage(t *testing.T) {
	t.Run("should parse message correctly", func(t *testing.T) {
		sender := "damilola"
		content := "Emphasis on the DEEP"
		timeStamp := "21/05/2024, 18:52:11"

		messageContent := fmt.Sprintf("[%s] %s: %s", timeStamp, sender, content)
		message, err := fileProcessor.ParseMessage(messageContent)
		if err != nil {
			t.Errorf("Error parsing message: %v", err)
		}
		if message.Sender != sender {
			t.Errorf("Expected sender to be %s, got %s", sender, message.Sender)
		}
		if message.Content != content {
			t.Errorf("Expected content to be %s, got %s", content, message.Content)
		}

		expectedDatetime, _ := time.Parse("02/01/2006, 15:04:05", timeStamp)
		if !message.Timestamp.Equal(expectedDatetime) {
			t.Errorf("Expected timestamp to be %s, got %s", expectedDatetime, message.Timestamp)
		}
	})

	t.Run("Should parse correctly even with delimiter in message", func(t *testing.T) {
		sender := "damilola"
		content := "Emphasis: on the DEEP"
		timeStamp := "21/05/2024, 18:52:11"

		messageContent := fmt.Sprintf("[%s] %s: %s", timeStamp, sender, content)
		message, _ := fileProcessor.ParseMessage(messageContent)
		if message.Content != content {
			t.Errorf("Expected content to be %s, got %s", content, message.Content)
		}
	})
}

func TestReadFile(t *testing.T) {
	t.Run("should read file correctly", func(t *testing.T) {
		t.Run("Should gather messages", func(t *testing.T) {
			conversation, _ := fileProcessor.Readfile("./testdata/test-pure.txt")
			if conversation.Sender != "test" {
				t.Errorf("Expected sender to be test, got %s", conversation.Sender)
			}
			if len(conversation.Messages) != 5 {
				t.Errorf("Expected 5 messages, got %d", len(conversation.Messages))
			}
		})

		t.Run("Should handle line breaks", func(t *testing.T) {
			conversation, _ := fileProcessor.Readfile("./testdata/test-with-line-breaks.txt")
			if len(conversation.Messages) != 5 {
				t.Errorf("Expected 5 messages, got %d", len(conversation.Messages))
			}
			expectedMessage := "Hey,\nhow\nare\nyou?"
			if conversation.Messages[0].Content != expectedMessage {
				t.Errorf("Expected message to be %s, got %s", expectedMessage, conversation.Messages[0].Content)
			}
		})
	})
}
