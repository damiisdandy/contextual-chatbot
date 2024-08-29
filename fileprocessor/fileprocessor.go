package fileprocessor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ctxh "github.com/damiisdandy/contextual-chatbot/contexthandler"
)

const (
	WhatsappDataPrivacyInfo = "Messages and calls are end-to-end encrypted. No one outside of this chat, not even WhatsApp, can read or listen to them."
)

type FileProcessor struct {
	ConversationFile string
	// name of the user using the chatbot
	Reciever string
}

func NewFileProcessor(reciever, conversationsPath string) *FileProcessor {
	return &FileProcessor{
		ConversationFile: filepath.Join("./conversations", conversationsPath),
		Reciever:         reciever,
	}
}

func (fp *FileProcessor) ParseMessage(message string) (ctxh.Message, error) {
	// Get timestamp
	start := strings.Index(message, "[")
	end := strings.Index(message, "]")

	if start == -1 || end == -1 {
		return ctxh.Message{}, errors.New("Invalid message format, missing timestamp")
	}

	datetime := strings.Replace(message[start+1:end], ", ", " ", 1)
	timeStamp, err := time.Parse("02/01/2006 15:04:05", datetime)

	if err != nil {
		return ctxh.Message{}, errors.New("Invalid message format, invalid timestamp")
	}

	// Get Sender and Message
	senderAndMessage := message[end+2:]
	colonSeparator := strings.Index(senderAndMessage, ": ")
	if colonSeparator == -1 {
		return ctxh.Message{}, errors.New("Invalid message format, missing sender")
	}
	sender := senderAndMessage[:colonSeparator]
	content := senderAndMessage[colonSeparator+2:]

	return ctxh.Message{
		Timestamp: timeStamp,
		Sender:    sender,
		Content:   content,
	}, nil
}

func (fp *FileProcessor) Readfile(file string) ([]ctxh.Message, error) {
	if !strings.Contains(file, ".txt") {
		return nil, fmt.Errorf("File must be a Whatsapp conversation file")
	}
	readFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %s", err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	messages := []ctxh.Message{}
	lineNumber := 0

	for fileScanner.Scan() {
		line := fileScanner.Text()
		message, err := fp.ParseMessage(line)
		if err != nil {
			if lineNumber == 0 {
				return nil, fmt.Errorf("Error reading file: %s, are you sure it is a Whatsapp conversation file?", err)
			}
			// handle line breaks in message
			messages[lineNumber-1].Content += "\n" + line
			continue
		}
		lineNumber++
		messages = append(messages, message)
	}

	// remove Whatsapp data privacy info
	if strings.Contains(messages[0].Content, WhatsappDataPrivacyInfo) {
		messages = messages[1:]
	}

	return messages, nil
}
