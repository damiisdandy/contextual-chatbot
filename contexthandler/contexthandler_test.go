package contexthandler

import (
	"encoding/json"
	"testing"
)

const AIResponse = `
[
    {
      "sender": "Gabby",
      "content": "Yh I get you",
      "timestamp": "2024-08-30T17:06:00+01:00"
    },
    {
      "sender": "Receiver",
      "content": "That's just how I am at time",
      "timestamp": "2024-08-30T17:06:00+01:00"
    },
    {
      "sender": "Sender",
      "content": "But I love attention",
      "timestamp": "2024-08-30T17:06:00+01:00"
    },
    {
      "sender": "System",
      "content": "I can give that but not every time, no one can give every time",
      "timestamp": "2024-08-30T17:06:00+01:00"
    },
    {
      "sender": "Receiver",
      "content": "You just notice when I don't",
      "timestamp": "2024-08-30T17:07:00+01:00"
    },
    {
      "sender": "Sender",
      "content": "I can give that but not every time, no one can give every time",
      "timestamp": "2024-08-30T17:09:00+01:00"
    },
    {
      "sender": "Sender",
      "content": "I understand",
      "timestamp": "2024-08-30T17:09:00+01:00"
    },
    {
      "sender": "Receiver",
      "content": "Yep",
      "timestamp": "2024-08-30T17:11:00+01:00"
    },
    {
      "sender": "Receiver",
      "content": "Pillow princess",
      "timestamp": "2024-08-30T17:10:00+01:00"
    }
  ]
`

var Peer = "Gabby ❤️❤️"

func TestAddMessages(t *testing.T) {
	messages := []Message{}
	_ = json.Unmarshal([]byte(AIResponse), &messages)
	t.Run("should properly identify the sender and reciever based on AI response [screenshot]", func(t *testing.T) {
		var contextStore = NewContextStore(Peer)
		contextStore.AddMessages(messages, MessageSourceScreenshot)

		senderTestTable := []struct {
			sender         string
			expectedSender string
		}{
			{contextStore.Messages[0].Sender, "Sender"},
			{contextStore.Messages[1].Sender, "Receiver"},
			{contextStore.Messages[2].Sender, "Sender"},
			{contextStore.Messages[3].Sender, "Receiver"},
		}

		for _, senderTest := range senderTestTable {
			if senderTest.sender != senderTest.expectedSender {
				t.Errorf("Expected sender to be %s, got %s", senderTest.expectedSender, senderTest.sender)
			}
		}
	})

	t.Run("should properly identify the sender and reciever based on chat logs", func(t *testing.T) {
		// simulate parsed chat logs
		for i, message := range messages {
			if i%2 == 0 {
				message.Sender = "damilola"
			} else {
				message.Sender = Peer
			}
		}

		var contextStore = NewContextStore(Peer)
		contextStore.AddMessages(messages, MessageSourceLogs)

		senderTestTable := []struct {
			sender         string
			expectedSender string
		}{
			{contextStore.Messages[0].Sender, "Sender"},
			{contextStore.Messages[1].Sender, "Receiver"},
			{contextStore.Messages[2].Sender, "Sender"},
			{contextStore.Messages[3].Sender, "Receiver"},
		}

		for _, senderTest := range senderTestTable {
			if senderTest.sender != senderTest.expectedSender {
				t.Errorf("Expected sender to be %s, got %s", senderTest.expectedSender, senderTest.sender)
			}
		}
	})

}
