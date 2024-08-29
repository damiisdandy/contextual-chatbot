package contexthandler

import "time"

type Message struct {
	Timestamp time.Time
	Sender    string
	Content   string
}

type Conversation struct {
	Sender   string
	Messages []Message
}
