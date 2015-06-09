package xingapi

type ConversationsList struct {
	Total         string          `json:"total"`
	Conversations []*Conversation `json:"items"`
}
