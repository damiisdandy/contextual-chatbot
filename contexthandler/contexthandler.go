package contexthandler

import (
	"fmt"
	"sort"
	"time"
)

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
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

func (cs *ContextStore) AddMessages(messages []Message) {
	for _, message := range messages {
		if message.Sender == cs.Peer {
			message.Sender = "Sender"
		} else {
			message.Sender = "Receiver"
		}
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
		chatLog += fmt.Sprintf("[%s] %s: %s\n", message.Timestamp, message.Sender, message.Content)
	}
	pastQuestions := ""
	for _, question := range cs.PastQuestions {
		pastQuestions += fmt.Sprintf("- %s\n", question)
	}
	if len(cs.PastQuestions) == 0 {
		pastQuestions = "None for now"
	}

	rootPrompt := fmt.Sprintf(`
		You are a professional relationship advisor, I am going to ask you questions about my relationship with %[1]s.

		You will responsed based on the following data:
		1. Past chat logs between I and %[1]s, where %[1]s is the Sender and I am the Receiver (chat logs are below)
		2. Past questions I've asked you (past questions are below)

		Chat logs:
		%[2]s

		Past Questions:
		%[3]s

		My Current Question:
		%[4]s

		Other things to consider:
		- Try to friendly, funny and engaging
		- Reference the past chat logs and past questions
		- Be concise
		- Do not repeat yourself
		- Mention the sender by their name (%[1]s)
	`, cs.Peer, chatLog, pastQuestions, question)

	// add new question to the context store
	cs.AddQuestion(question)
	return rootPrompt
}
