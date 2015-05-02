package xingapi

import (
	"bufio"
	"os"
)

// The file name for the file to store the credentials in
const credentialsFileName = "credentialsFile.txt"

// Credentials are used to authenticate against the API via OAuth
type Credentials struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}

// CredentialStore stores credentials
type CredentialStore struct{}

// SaveCredentials saves the given credentials
func (store *CredentialStore) SaveCredentials(credentials Credentials) error {
	var marshaler CredentialsMarshaler
	marshaler = JSONMarshaler{}
	file, _ := os.Create(credentialsFileName)
	writer := bufio.NewWriter(file)
	err := marshaler.MarshalCredentials(writer, credentials)
	writer.Flush()
	return err
}

// Credentials returns the stored credentials or an error if no credentials have been stored
func (store *CredentialStore) Credentials() (Credentials, error) {
	var unmarshaler CredentialsUnmarshaler
	unmarshaler = JSONMarshaler{}
	file, err := os.Open(credentialsFileName)
	reader := bufio.NewReader(file)
	credentials, err := unmarshaler.UnmarshalCredentials(reader)
	return credentials, err
}
