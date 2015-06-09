package xingapi

type ConversationsList struct {
	Total         int             `json:"total"`
	Conversations []*Conversation `json:"items"`
}
