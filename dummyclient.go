package xingapi

// DummyClient is used as a mock for tests the need a Client
type DummyClient struct {
	Client
	DummyUsers []User
}

// ContactsList is a fake ContactsList implementation
func (client *DummyClient) ContactsList(userID string, limit int, offset int, handler ContactsHandler) {
	list := new(ContactsList)
	list.UserIDs = []string{"userId 1", "userId 2"}
	list.Total = 2
	handler(*list, nil)
}

// User is a fake User implementation
func (client *DummyClient) User(contactUserID string, handler UserHandler) {
	handler(client.DummyUsers[0], nil)
}
