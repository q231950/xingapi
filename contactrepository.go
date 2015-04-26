// contactrespository.go

package xingapi

import (
	"fmt"
	"github.com/str1ngs/ansi/color"

	"sync"
)

type ContactRepository struct {
	client Client
}

func NewContactRepository(client Client) *ContactRepository {
	return &ContactRepository{client}
}

func (repo *ContactRepository) Contacts(userId string, contactsHandler func(list []*User, err error)) {
	repo.client.ContactsList(userId, 0, 0, func(list ContactsList, err error) {
		if err == nil {
			color.Printf("", fmt.Sprintf("-----------------------------------\n%d Contacts\n", list.Total))
			if 0 < list.Total {
				limit := 50
				request := UsersRequest{userId, limit, 0, list.Total, contactsHandler}
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

	newRequest := UsersRequest{request.UserId,
		limit,
		request.Offset,
		request.Total,
		request.Completion}
	repo.loadUsers([]*User{}, newRequest)
}

func (repo *ContactRepository) loadUsers(users []*User, request UsersRequest) {
	repo.client.ContactsList(request.UserId, request.Limit, request.Offset, func(list ContactsList, err error) {
		if err == nil {
			repo.loadUserDetails(list, func(loadedUsers []*User, err error) {
				users = append(users, loadedUsers...)
				if !request.IsFinal() {
					newRequest := UsersRequest{request.UserId, request.Limit, request.Offset + len(list.UserIds), request.Total, request.Completion}
					repo.loadUsers(users, newRequest)
				} else {
					// finished final request without errors
					request.Completion(users, nil)
				}
			})
		} else {
			request.Completion(nil, err)
		}
	})
}

func (repo *ContactRepository) loadUserDetails(list ContactsList, loadedUsers func(userList []*User, err error)) {
	users := []*User{}
	//	var err error
	var waitGroup sync.WaitGroup
	for _, contactUserId := range list.UserIds {
		waitGroup.Add(1)
		go repo.client.User(contactUserId, func(user User, cerr error) {
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
