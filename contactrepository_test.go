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
	firstDummy := new(DummyUser)
	firstDummy.UserID = "userId 1"
	dummyUsers[0] = firstDummy

	secondDummy := new(DummyUser)
	secondDummy.UserID = "userId 2"
	dummyUsers[1] = secondDummy
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
