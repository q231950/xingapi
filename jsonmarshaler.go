/*
Package xingapi contains the JSONMarshaler that allows marshaling and unmarshaling all entities to and from JSON
*/
package xingapi

import (
	"encoding/json"
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

/*
JSONMarshaler is a concrete implementation of
- UsersMarshaler/UsersUnmarshaler
- CredentialsMarshaler/CredentialsUnmarshaler
- ContactsListUnmarshaler
- UserUnmarshaler
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
		list.Total = jsonlist.JSONContactsUserIdList.Total
		list.UserIDs = jsonlist.JSONContactsUserIdList.UserIds()
	}
	return *list, err
}

// UnmarshalUser concrete implementation
func (JSONMarshaler) UnmarshalUser(reader io.Reader) (User, error) {
	var marshaler UsersUnmarshaler
	marshaler = JSONMarshaler{}
	users, err := marshaler.UnmarshalUsers(reader)
	return *users.Users[0], err
}
