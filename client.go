package xingapi

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
	"sync"
)

// A UserHandler is used as parameter to methods that fetch a User. Either the User or error might be nil.
type UserHandler func(User, error)

// The ContactsHandler can be passed as parameter to methods fetching contact lists. Either ContactsList or error might be nil.
type ContactsHandler func(ContactsList, error)

// The Client interface describes the available methods that wrap around the XING API.
type Client interface {
	User(id string, handler UserHandler)
	ContactsList(userID string, limit int, offset int, handler ContactsHandler)
	Me(handler UserHandler)
	Messages(userID string, handler func(err error))
}

/*
The XINGClient conforms to the Client interface and handles oAuth authentication for the XING API.
It manages all communication with the API's endpoints.
*/
type XINGClient struct {
	OAuthConsumer OAuthConsumer
}

// Me fetches the logged in user
func (client *XINGClient) Me(handler UserHandler) {
	var me User
	consumer := new(OAuthConsumer)
	client.OAuthConsumer = *consumer
	client.OAuthConsumer.Get("/v1/users/me", url.Values{}, func(reader io.Reader, err error) {
		var unmarshaler UsersUnmarshaler
		unmarshaler = JSONMarshaler{}

		users, jsonError := unmarshaler.UnmarshalUsers(reader)
		if jsonError != nil {
			err = jsonError
		}

		me = *users.Users[0]
		handler(me, err)
	})
}

// ContactsList fetches the contact list of the logged in user from the offset with a batch count of limit
func (client *XINGClient) ContactsList(userID string, limit int, offset int, handler ContactsHandler) {
	fmt.Println("Loading contacts list from server...")
	consumer := new(OAuthConsumer)
	client.OAuthConsumer = *consumer
	v := url.Values{}
	v.Set("limit", strconv.Itoa(limit))
	v.Set("offset", strconv.Itoa(offset))
	v.Set("order_by", "last_name")
	client.OAuthConsumer.Get("/v1/users/"+userID+"/contacts", v, func(reader io.Reader, err error) {

		var unmarshaler ContactsListUnmarshaler
		unmarshaler = JSONMarshaler{}
		list, jsonError := unmarshaler.UnmarshalContactsList(reader)
		if jsonError != nil {
			fmt.Println("Error while loading contacts list")
			err = jsonError
		}
		handler(list, err)
	})
}

// User fetches the User for the given user id
func (client *XINGClient) User(id string, handler UserHandler) {
	fmt.Println("client gets user <" + id + ">")
	consumer := new(OAuthConsumer)
	client.OAuthConsumer = *consumer
	client.OAuthConsumer.Get("/v1/users/"+id, url.Values{}, func(reader io.Reader, err error) {
		var unmarshaler UserUnmarshaler
		unmarshaler = JSONMarshaler{}
		user, jsonError := unmarshaler.UnmarshalUser(reader)
		if jsonError != nil {
			err = jsonError
			PrintError(err)
		}
		handler(user, err)
	})
}

// Users fetches the users according to the given user IDs
func (client *XINGClient) Users(userIDs []string) []User {
	users := []User{}
	var waitGroup sync.WaitGroup
	for _, userID := range userIDs {
		waitGroup.Add(1)
		go client.User(userID, func(user User, err error) {
			if err == nil {
				PrintUserOneLine(user)
				users = append(users, user)
			}
			waitGroup.Done()
		})
	}

	waitGroup.Wait()
	return users
}

// Messages fetches the conversations of the user with the given id and prints out raw json
func (client *XINGClient) Messages(userID string, handler func(err error)) {
	client.OAuthConsumer.Get("/v1/users/"+userID+"/conversations", url.Values{}, func(reader io.Reader, err error) {
		fmt.Println("begin messages")
		var unmarshaler = ConversationsMarshaler{}
		conversationsInfo, JSONError := unmarshaler.UnmarshalConversationList(reader)
		if JSONError == nil {
			for _, Conversation := range conversationsInfo.ConversationsList.Conversations {
				fmt.Println(Conversation)
				userIDs := []string{}
				userID_0 := Conversation.Participants[0].UserID
				userID_1 := Conversation.Participants[1].UserID
				userIDs = append(userIDs, userID_0)
				userIDs = append(userIDs, userID_1)
				client.Users(userIDs)
				fmt.Println("-------")
			}
		} else {
			PrintError(JSONError)
		}
		handler(JSONError)
	})
}
