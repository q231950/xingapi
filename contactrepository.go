package xingapi

import "sync"

// ContactRepository represents a repository to retrieve and store contacts
type ContactRepository struct {
	client Client
}

// NewContactRepository creates a new contact repository given a API client
func NewContactRepository(client Client) *ContactRepository {
	return &ContactRepository{client}
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
					newRequest := UsersRequest{request.UserId, request.Limit, request.Offset + len(list.UserIDs), request.Total, request.Completion}
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
