package xingapi

import (
	"encoding/json"
	"io"
)

type ConversationsMarshaler struct{}

func (ConversationsMarshaler) UnmarshalConversationList(reader io.Reader) (ConversationsInfo, error) {
	decoder := json.NewDecoder(reader)
	var conversationsList ConversationsInfo
	err := decoder.Decode(&conversationsList)
	return conversationsList, err
}
