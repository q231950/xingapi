package xingapi

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/boltdb/bolt"
)

// ContactRepository represents a repository to retrieve and store contacts
type ContactRepository struct {
	client Client
	db     *bolt.DB
}

// NewContactRepository creates a new contact repository given a API client
func NewContactRepository(client Client) *ContactRepository {
	db, err := bolt.Open("contacts.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	return &ContactRepository{client, db}
}

// Contacts fetches all contacts of the user with the given user ID
func (repo *ContactRepository) Contacts(userID string, contactsHandler func(list []*User, err error)) {

	repo.client.ContactsList(userID, 0, 0, func(list ContactsList, err error) {
		if err == nil {
			if 0 < list.Total {
				limit := 50
				request := UsersRequest{userID, limit, 0, list.Total, contactsHandler}
				repo.requestLoadUsers(request)
			} else {
				contactsHandler([]*User{}, nil)
			}
		} else {
			contactsHandler(nil, err)
		}
	})
}

func (repo *ContactRepository) requestLoadUsers(request UsersRequest) {

	limit := request.Limit
	if request.Offset+request.Limit > request.Total {
		limit = request.Limit - (request.Offset + request.Limit - request.Total)
	}

	newRequest := UsersRequest{request.UserID,
		limit,
		request.Offset,
		request.Total,
		request.Completion}
	repo.loadUsers([]*User{}, newRequest)
}

func (repo *ContactRepository) loadUsers(users []*User, request UsersRequest) {
	repo.client.ContactsList(request.UserID, request.Limit, request.Offset, func(list ContactsList, err error) {
		if err == nil {
			repo.loadUserDetails(list, func(loadedUsers []*User, err error) {
				users = append(users, loadedUsers...)
				if !request.IsFinal() {
					newRequest := UsersRequest{request.UserID, request.Limit, request.Offset + len(list.UserIDs), request.Total, request.Completion}
					repo.loadUsers(users, newRequest)
				} else {
					// finished final request without errors
					repo.storeUsers(users)
					request.Completion(users, nil)
				}
			})
		} else {
			request.Completion(nil, err)
		}
	})
}

var world = []byte("nameBucket")

func (repo *ContactRepository) storeUsers(users []*User) {
	for _, user := range users {

		var marshaler UserMarshaler
		marshaler = JSONMarshaler{}
		// writer := bufio.NewWriter()
		bytes, err := marshaler.MarshalUser(*user)
		if err == nil {
			key := []byte((*user).DisplayName())
			value := bytes

			err = repo.db.Update(func(tx *bolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists(world)
				if err != nil {
					return err
				}

				err = bucket.Put(key, value)
				if err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Contact gets the user with the given name
func (repo *ContactRepository) Contact(name string, done func()) {

	key := []byte(name)
	// retrieve the data
	err := repo.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(world)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", world)
		}

		val := bucket.Get(key)
		fmt.Println(string(val))

		c := bucket.Cursor()

		prefix := []byte(name)
		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		done()
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (repo *ContactRepository) loadUserDetails(list ContactsList, loadedUsers func(userList []*User, err error)) {
	users := []*User{}
	//	var err error
	var waitGroup sync.WaitGroup
	for _, contactUserID := range list.UserIDs {
		waitGroup.Add(1)
		go repo.client.User(contactUserID, func(user User, cerr error) {
			if cerr == nil {
				users = append(users, &user)
			} else {
				PrintError(cerr)
				//				err = cerr
			}
			defer waitGroup.Done()
		})
	}
	waitGroup.Wait()
	loadedUsers(users, nil)
}
