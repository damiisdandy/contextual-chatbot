package screenshotprocessor

import "testing"

var screenShotProcess = ScreenshotProcessor{}

func TestParseJSONString(t *testing.T) {
	mockAIResponse := `[{"sender":"Sender","content":"Call Him","timestamp":"2024-08-29T22:12:00+01:00"},{"sender":"Sender","content":"Voice call\nNo answer","timestamp":"2024-08-29T03:24:00+01:00"},{"sender":"Sender","content":"Voice call\n15 sec","timestamp":"2024-08-29T03:26:00+01:00"},{"sender":"Receiver","content":"What's happening to boiski event","timestamp":"2024-08-29T13:29:00+01:00"},{"sender":"Receiver","content":"i have no idea, i want to text the\nguy now","timestamp":"2024-08-29T13:29:00+01:00"},{"sender":"Receiver","content":"i checked the site everything is fine\nw it","timestamp":"2024-08-29T13:30:00+01:00"},{"sender":"Receiver","content":"Jetron Ticket | BIOSKY'S PLAYHOUSE\nIf you came for the last edition you'd know\nbest not to miss this one . September 1st we...\n\nwww.jetronticket.com\n\nhttps://www.jetronticket.com/\nevents/bioskys-playhouse-1","timestamp":"2024-08-29T13:30:00+01:00"}]`
	t.Run("should parse JSON string correctly", func(t *testing.T) {
		messages, err := screenShotProcess.ParseJSONString(mockAIResponse)
		if err != nil {
			t.Fatalf("Error parsing JSON string: %s", err)
		}
		if len(messages) != 7 {
			t.Errorf("Expected 7 messages, got %d", len(messages))
		}
		expectedMessage := "Call Him"
		if messages[0].Content != expectedMessage {
			t.Errorf("Expected message to be %s, got %s", expectedMessage, messages[0].Content)
		}
	})

	t.Run("should parse correctly even if AI adds statement", func(t *testing.T) {
		_, err := screenShotProcess.ParseJSONString("Here is the provided result:\n" + mockAIResponse)
		if err != nil {
			t.Fatalf("Error parsing JSON string: %s", err)
		}
	})

	t.Run("remove duplicate because AI sees a reply as a new message", func(t *testing.T) {
		mockAIResponse = `[{"sender":"Sender","content":"Call Him","timestamp":"2024-08-29T22:12:00+01:00"},{"sender":"Reciever","content":"Call Him","timestamp":"2024-08-29T22:12:00+01:00"},{"sender":"Reciever","content":"Okay","timestamp":"2024-08-29T22:12:00+01:00"}]`
		messages, _ := screenShotProcess.ParseJSONString(mockAIResponse)
		if len(messages) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(messages))
		}
	})
}
