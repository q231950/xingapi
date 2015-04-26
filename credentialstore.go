// credentialstore.go

package xingapi

import (
	"os"
	"bufio"
	)

const credentialsFileName = "credentialsFile.txt"

type Credentials struct {
	Token string `json:"token"`
	Secret string `json:"secret"`
}

type CredentialStore struct {}

func (store *CredentialStore)SaveCredentials(credentials Credentials) error {
	var marshaler CredentialsMarshaler
	marshaler = JSONMarshaler{}
	file, _ := os.Create(credentialsFileName)
	writer := bufio.NewWriter(file)
	err := marshaler.MarshalCredentials(writer, credentials)
	writer.Flush()
	return err
}

func (store *CredentialStore)Credentials() (Credentials, error) {
	var unmarshaler CredentialsUnmarshaler
	unmarshaler = JSONMarshaler{}
	file, err := os.Open(credentialsFileName)
	reader := bufio.NewReader(file)
	credentials, err := unmarshaler.UnmarshalCredentials(reader)
	return credentials, err
}
