package contexthandler

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const MessageSourceScreenshot = "screenshot"
const MessageSourceLogs = "logs"

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Source    string    `json:"source"`
}

type ContextStore struct {
	Messages []Message
	// who we are talking to
	Peer string
	// questions we've asked the AI
	PastQuestions []string
}

func NewContextStore(peer string) *ContextStore {
	return &ContextStore{
		Messages:      []Message{},
		Peer:          peer,
		PastQuestions: []string{},
	}
}

func (cs *ContextStore) AddMessages(messages []Message, source string) {
	for _, message := range messages {
		sender := strings.TrimSpace(message.Sender)
		if sender != "Sender" && sender != "Receiver" {
			// e.g Gabby ❤️❤️ vs Gabby
			if strings.Contains(cs.Peer, sender) {
				message.Sender = "Sender"
			} else {
				message.Sender = "Receiver"
			}
		}
		message.Source = source
		cs.Messages = append(cs.Messages, message)
	}

	sort.Slice(cs.Messages, func(i, j int) bool {
		return cs.Messages[j].Timestamp.After(cs.Messages[i].Timestamp)
	})
}

func (cs *ContextStore) AddQuestion(question string) {
	cs.PastQuestions = append(cs.PastQuestions, question)
}

func (cs *ContextStore) GeneratePromp(question string) string {
	chatLog := ""
	for _, message := range cs.Messages {
		chatLog += fmt.Sprintf("[%s] [source: %s] %s: %s\n", message.Timestamp, message.Source, message.Sender, message.Content)
	}
	pastQuestions := ""
	for _, question := range cs.PastQuestions {
		pastQuestions += fmt.Sprintf("- %s\n", question)
	}
	if len(cs.PastQuestions) == 0 {
		pastQuestions = "None for now"
	}

	rootPrompt := fmt.Sprintf(`
		You are a dating assistant, I am going to ask you questions about my relationship with %[1]s.

		You will responsed based on the following data:
		1. Past chat logs between I and %[1]s, where %[1]s is the Sender and I am the Receiver (chat logs are below)
		2. Past questions I've asked you (past questions are below)
		3. Each log has a source (%[5]s or %[6]s)

		<chat-logs>
		%[2]s
		</chat-logs>

		<past-questions>
		%[3]s
		</past-questions>

		My Current Question:
		%[4]s

		Other things to consider:
		- Reference the past chat logs and past questions.
		- when I mention screenshot, focus on the chat logs that have the source "%[6]s".
		- keep track of the order of each screenshots and chat logs based on their timestamps.

		- Be short and concise, reply like we are texting (keep your response short and to the point).
		- your response should be less than 50 words.

		- Do not mention the fact that you based your response on chat logs and past chats.
		- Do not use the word "our" or "Based on" in your response.

		- Mention the sender by their name (%[1]s).
		- You are to reply as a third-party dating assistant analysing me and %[1]s's relationship.

		- Give example of chats that drew your conclusion, also mention its source.
	`, cs.Peer, chatLog, pastQuestions, question, MessageSourceLogs, MessageSourceScreenshot)

	// add new question to the context store
	cs.AddQuestion(question)
	return rootPrompt
}
