package xingapi

import (
	"encoding/json"
	"io"
)

type JSONMessageMarshaler struct{}

func (JSONMessageMarshaler) UnmarshalMessage(reader io.Reader) (Message, error) {
	decoder := json.NewDecoder(reader)
	var message Message
	err := decoder.Decode(&message)
	return message, err
}
