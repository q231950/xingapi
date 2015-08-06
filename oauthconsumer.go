package xingapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ResponseHandler handles HTTP response bodies and a potential error
type ResponseHandler func(io.Reader, error)

// AuthenticateHandler will be called upon successful authentication
type AuthenticateHandler func()

/*
OAuthConsumer connects to the API in a authenticated manner. It wraps GET calls
to the API by attaching OAuth information to them.
*/
type OAuthConsumer struct {
	// Authenticated indicates whether or not the consumer is currently authenticated
	Authenticated      bool
	oauthAuthenticator *OAuthAuthenticator

	// AuthenticateHandlers store handlers to be executed once authentication was successful
	AuthenticateHandlers []AuthenticateHandler
}

func (consumer *OAuthConsumer) authenticate(handler AuthenticateHandler) {

	if consumer.Authenticated {
		handler()
	} else {
		credentialStore := new(CredentialStore)
		storedCredentials, localeCredentialsError := credentialStore.Credentials()

		if localeCredentialsError == nil {
			consumer.oauthAuthenticator = new(OAuthAuthenticator)
			consumer.oauthAuthenticator.AuthenticateUsingStoredCredentials(storedCredentials, func(err error) {
				if err == nil {
					consumer.handleAuthentication()
				}
			})
		} else {
			consumer.requestCredentials(handler)
		}
	}
}

func (consumer *OAuthConsumer) requestCredentials(handler func()) {
	consumer.oauthAuthenticator = new(OAuthAuthenticator)
	consumer.oauthAuthenticator.Authenticate(func(err error) {
		if err == nil {
			handler()
		}
	})
}

// Get performs get requests to the API while using OAuth authentication.
func (consumer *OAuthConsumer) Get(path string, parameters url.Values, handler ResponseHandler) {
	consumer.addAuthenticationHandler(func() {
		httpClient := new(http.Client)
		url := "https://api.xing.com" + path
		credentials := consumer.oauthAuthenticator.OAuthCredentials
		resp, err := consumer.oauthAuthenticator.Client.Get(httpClient, &credentials, url, parameters)
		if resp.StatusCode == 200 {
			handler(resp.Body, err)
		} else {
			PrintCommand(fmt.Sprintf("GET %s\n", path))
			//contents, _ := ioutil.ReadAll(resp.Body)
			//fmt.Printf("%s\n", string(contents))

			var apiError APIError
			responseJSON, _ := ioutil.ReadAll(resp.Body)
			err := json.Unmarshal(responseJSON, &apiError)
			if err == nil {
				fmt.Println(apiError)
			}
		}
		defer resp.Body.Close()
	})

	if !consumer.Authenticated {
		consumer.authenticate(func() {
			consumer.handleAuthentication()
		})
	}
}

func (consumer *OAuthConsumer) addAuthenticationHandler(handler AuthenticateHandler) {
	if consumer.AuthenticateHandlers == nil {
		consumer.AuthenticateHandlers = []AuthenticateHandler{}
	}

	if consumer.Authenticated {
		handler()
	} else {
		consumer.AuthenticateHandlers = append(consumer.AuthenticateHandlers, handler)
	}
}

func (consumer *OAuthConsumer) handleAuthentication() {
	consumer.Authenticated = true
	for _, authenticateHandler := range consumer.AuthenticateHandlers {
		authenticateHandler()
	}
}
