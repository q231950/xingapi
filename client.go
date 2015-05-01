package xingapi

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
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
	Messages(userId string, handler func(err error))
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
			err = jsonError
		}
		handler(list, err)
	})
}

// User fetches the User for the given user id
func (client *XINGClient) User(id string, handler UserHandler) {
	consumer := new(OAuthConsumer)
	client.OAuthConsumer = *consumer
	client.OAuthConsumer.Get("/v1/users/"+id, url.Values{}, func(reader io.Reader, err error) {
		var unmarshaler UserUnmarshaler
		unmarshaler = JSONMarshaler{}
		user, jsonError := unmarshaler.UnmarshalUser(reader)
		if jsonError != nil {
			err = jsonError
		}
		handler(user, err)
	})
}

// Messages fetches the conversations of the user with the given id and prints out raw json
func (client *XINGClient) Messages(userId string, handler func(err error)) {
	client.OAuthConsumer.Get("/v1/users/"+userId+"/conversations", url.Values{}, func(reader io.Reader, err error) {
		robots, readError := ioutil.ReadAll(reader)

		println(fmt.Sprintf("%s", robots))

		handler(readError)
	})
}
