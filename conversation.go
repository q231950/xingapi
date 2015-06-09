package xingapi

import (
	"strconv"
	"time"
)

type Conversation struct {
	Subject          string `json:"subject"`
	MessageCount     int    `json:"message_count"`
	ChangeTimeString string `json:"updated_at"`
	Readonly         bool   `json:"read_only"`
	Participants     []struct {
		UserID string `json:"id"`
	} `json:"participants"`
}

func (c Conversation) String() string {

	time, err := c.ChangeTime()
	timeString := time.String()
	if err != nil {
		timeString = "Unknown date"
	}
	return timeString + ":" + c.Subject + " between " + c.Participants[0].UserID + " and " + c.Participants[1].UserID + " with " + strconv.Itoa(c.MessageCount) + " messages."
}

func (c Conversation) ChangeTime() (time.Time, error) {
	return time.Parse(time.RFC3339Nano, c.ChangeTimeString)
}
