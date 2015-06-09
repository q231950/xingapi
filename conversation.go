package xingapi

import "strconv"

type Conversation struct {
	Subject      string `json:"subject"`
	MessageCount int    `json:"message_count"`
}

func (cl Conversation) String() string {
	return "Conversation:" + cl.Subject + " with " + strconv.Itoa(cl.MessageCount) + " messages."
}
