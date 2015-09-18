package xingapi

import (
	"encoding/json"
	"fmt"
	"io"
)

// UsersMarshaler interface defines Users marshaler
type UsersMarshaler interface {
	MarshalUsers(writer io.Writer, users Users) error
}

// UsersUnmarshaler interface defines Users unmarshaler
type UsersUnmarshaler interface {
	UnmarshalUsers(reader io.Reader) (Users, error)
}

// CredentialsMarshaler defines Credentials marshaler
type CredentialsMarshaler interface {
	MarshalCredentials(writer io.Writer, credentials Credentials) error
}

// CredentialsUnmarshaler defines Credentials unmarshaler
type CredentialsUnmarshaler interface {
	UnmarshalCredentials(reader io.Reader) (Credentials, error)
}

// ContactsListUnmarshaler defines ContactsList marshaler
type ContactsListUnmarshaler interface {
	UnmarshalContactsList(reader io.Reader) (ContactsList, error)
}

// UserUnmarshaler defines User marshaler
type UserUnmarshaler interface {
	UnmarshalUser(reader io.Reader) (User, error)
}

// UserMarshaler defines marshaler for marshaling user to JSON
type UserMarshaler interface {
	MarshalUser(user User) (bytes []byte, err error)
}

/*
The JSONMarshaler is a concrete implementation of
UsersMarshaler/UsersUnmarshaler,
CredentialsMarshaler/CredentialsUnmarshaler,
ContactsListUnmarshaler and UserUnmarshaler.
The JSONMarshaler that allows marshaling and unmarshaling all entities to and from JSON.
*/
type JSONMarshaler struct{}

// MarshalUsers concrete implementation
func (JSONMarshaler) MarshalUsers(writer io.Writer, users Users) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(users)
}

// UnmarshalUsers concrete implementation
func (JSONMarshaler) UnmarshalUsers(reader io.Reader) (Users, error) {
	decoder := json.NewDecoder(reader)
	var users Users
	err := decoder.Decode(&users)
	return users, err
}

// MarshalCredentials concrete implementation
func (JSONMarshaler) MarshalCredentials(writer io.Writer, credentials Credentials) error {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(credentials)
	return err
}

// UnmarshalCredentials concrete implementation
func (JSONMarshaler) UnmarshalCredentials(reader io.Reader) (Credentials, error) {
	decoder := json.NewDecoder(reader)
	var credentials Credentials
	err := decoder.Decode(&credentials)
	return credentials, err
}

// UnmarshalContactsList concrete implementation
func (JSONMarshaler) UnmarshalContactsList(reader io.Reader) (ContactsList, error) {

	decoder := json.NewDecoder(reader)
	var jsonlist JSONContactsList
	err := decoder.Decode(&jsonlist)
	list := new(ContactsList)

	if err == nil {
		list.Total = jsonlist.JSONContactsUserIDList.Total
		list.UserIDs = jsonlist.JSONContactsUserIDList.UserIDs()
	}
	return *list, err
}

// UnmarshalUser concrete implementation
func (JSONMarshaler) UnmarshalUser(reader io.Reader) (User, error) {
	var marshaler UsersUnmarshaler
	marshaler = JSONMarshaler{}
	users, err := marshaler.UnmarshalUsers(reader)
	fmt.Printf("%d", len(users.Users))
	return users.Users[0], err
}

// MarshalUser concrete implementation
func (JSONMarshaler) MarshalUser(user User) ([]byte, error) {
	return json.Marshal(user)
}
