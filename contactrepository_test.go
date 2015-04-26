// contactsrepository_test.go

package xingapi

import (
	"strconv"
	"sync"
	"testing"
)

func TestGetContacts(t *testing.T) {
	client := new(DummyClient)
	dummyUsers := make([]User, 2)
	dummyUsers[0] = new(DummyUser)
	dummyUsers[1] = new(DummyUser)
	client.DummyUsers = dummyUsers
	repository := NewContactRepository(client)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	repository.Contacts("some user id", func(list []*User, err error) {
		if len(list) != 2 {
			t.Error("Expected '2' but got '" + strconv.Itoa(len(list)) + "'")
		}
		waitGroup.Done()
	})
	waitGroup.Wait()
}
