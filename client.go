// client.go
package xingapi

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
)

type UserHandler func(User, error)
type ContactsHandler func(ContactsList, error)

type Client interface {
	User(id string, handler UserHandler)
	ContactsList(userID string, limit int, offset int, handler ContactsHandler)
	Me(handler UserHandler)
	Messages(userId string, handler func(err error))
}

type XINGClient struct {
	OAuthConsumer OAuthConsumer
}

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

// GET /v1/users/:user_id/conversations
func (client *XINGClient) Messages(userId string, handler func(err error)) {
	client.OAuthConsumer.Get("/v1/users/"+userId+"/conversations", url.Values{}, func(reader io.Reader, err error) {
		robots, readError := ioutil.ReadAll(reader)

		println(fmt.Sprintf("%s", robots))

		handler(readError)
	})
}
