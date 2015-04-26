// jsonmarshaler.go
package xingapi

import (
	"encoding/json"
	"io"
)

type UsersMarshaler interface {
	MarshalUsers(writer io.Writer, users Users) error
}

type UsersUnmarshaler interface {
	UnmarshalUsers(reader io.Reader) (Users, error)
}

type CredentialsMarshaler interface {
	MarshalCredentials(writer io.Writer, credentials Credentials) error
}

type CredentialsUnmarshaler interface {
	UnmarshalCredentials(reader io.Reader) (Credentials, error)
}

type ContactsListUnmarshaler interface {
	UnmarshalContactsList(reader io.Reader) (ContactsList, error)
}

type UserUnmarshaler interface {
	UnmarshalUser(reader io.Reader) (User, error)
}

type JSONMarshaler struct{}

func (JSONMarshaler) MarshalUsers(writer io.Writer, users Users) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(users)
}

func (JSONMarshaler) UnmarshalUsers(reader io.Reader) (Users, error) {
	decoder := json.NewDecoder(reader)
	var users Users
	err := decoder.Decode(&users)
	return users, err
}

func (JSONMarshaler) MarshalCredentials(writer io.Writer, credentials Credentials) error {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(credentials)
	return err
}

func (JSONMarshaler) UnmarshalCredentials(reader io.Reader) (Credentials, error) {
	decoder := json.NewDecoder(reader)
	var credentials Credentials
	err := decoder.Decode(&credentials)
	return credentials, err
}

func (JSONMarshaler) UnmarshalContactsList(reader io.Reader) (ContactsList, error) {

	decoder := json.NewDecoder(reader)
	var jsonlist JSONContactsList
	err := decoder.Decode(&jsonlist)
	list := new(ContactsList)

	if err == nil {
		list.Total = jsonlist.JSONContactsUserIdList.Total
		list.UserIds = jsonlist.JSONContactsUserIdList.UserIds()
	}
	return *list, err
}

func (JSONMarshaler) UnmarshalUser(reader io.Reader) (User, error) {
	var marshaler UsersUnmarshaler
	marshaler = JSONMarshaler{}
	users, err := marshaler.UnmarshalUsers(reader)
	return *users.Users[0], err
}
