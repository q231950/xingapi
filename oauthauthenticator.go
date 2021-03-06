package xingapi

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/etix/stoppableListener"
	"github.com/garyburd/go-oauth/oauth"
)

/*
AuthenticationHandler is the function to be invoked once authentication
against the API has completed
*/
type AuthenticationHandler func(err error)

// OAuthAuthenticator handles OAuth authentication
type OAuthAuthenticator struct {
	Client               oauth.Client
	TemporaryCredentials oauth.Credentials
	OAuthCredentials     oauth.Credentials
	listener             *stoppableListener.StoppableListener
	authenticateHandler  AuthenticationHandler
}

/*
Authenticate starts the local server that waits for the OAuth callback once
the user signs in in the browser
*/
func (authenticator *OAuthAuthenticator) Authenticate(handler AuthenticationHandler) {

	authenticator.authenticateHandler = handler
	err := authenticator.getRequestToken()

	if err == nil {
		listener, _ := net.Listen("tcp", "127.0.0.1:8080")
		authenticator.listener = stoppableListener.Handle(listener)

		k := make(chan os.Signal, 1)
		signal.Notify(k, os.Interrupt)
		go func() {
			<-k
			authenticator.listener.Stop <- true
		}()

		http.HandleFunc("/", authenticator.onReceiveTemporaryVerifierAndToken)
		err = http.Serve(authenticator.listener, nil)
	}

	if err != nil {
		authenticator.authenticateHandler(err)
	}
}

// AuthenticateUsingStoredCredentials sets up the authenticator's client and OAuthCredentials. woot!?
func (authenticator *OAuthAuthenticator) AuthenticateUsingStoredCredentials(storedCredentials Credentials, handler AuthenticationHandler) {
	authenticator.Client.Credentials = authenticator.oAuthConsumerCredentials()
	authenticator.OAuthCredentials = oauth.Credentials{storedCredentials.Token, storedCredentials.Secret}
	handler(nil)
}

func (authenticator *OAuthAuthenticator) oAuthConsumerCredentials() oauth.Credentials {
	credentials := NewCredentials()
	return oauth.Credentials{credentials.Token, credentials.Secret}
}

func (authenticator *OAuthAuthenticator) oAuthClient() oauth.Client {
	client := new(oauth.Client)
	client.Credentials = authenticator.oAuthConsumerCredentials()
	client.TemporaryCredentialRequestURI = "https://api.xing.com/v1/request_token"
	client.ResourceOwnerAuthorizationURI = "https://api.xing.com/v1/authorize"
	client.TokenRequestURI = "https://api.xing.com/v1/access_token"
	client.SignatureMethod = oauth.PLAINTEXT
	return *client
}

func (authenticator *OAuthAuthenticator) getRequestToken() error {
	authenticator.Client = authenticator.oAuthClient()
	httpClient := new(http.Client)
	cred, err := authenticator.Client.RequestTemporaryCredentials(httpClient, "http://localhost:8080/", nil)
	if err == nil {
		authenticator.TemporaryCredentials = *cred
		tc := authenticator.TemporaryCredentials
		accessURL := authenticator.Client.AuthorizationURL(&tc, nil)
		PrintMessageWithParam("Please paste this url into your browser:", accessURL)
	}
	return err
}

func (authenticator *OAuthAuthenticator) onReceiveTemporaryVerifierAndToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	verifier, _ := r.Form["oauth_verifier"]
	token, _ := r.Form["oauth_token"]
	fmt.Fprintf(w, "<html><head></head><body>")
	if 0 < len(verifier) && 0 < len(token) {
		httpClient := new(http.Client)
		tc := authenticator.TemporaryCredentials
		credentials, _, err := authenticator.Client.RequestToken(httpClient, &tc, verifier[0])

		credentialStore := new(CredentialStore)
		credentialStore.SaveCredentials(Credentials{credentials.Token, credentials.Secret})

		authenticator.OAuthCredentials = *credentials
		authenticator.authenticateHandler(err)
		authenticator.listener.Stop <- true
		fmt.Fprintf(w, "<h3><span style=\"font-family:'Helvetica'; color:#777\">Success, you are authenticated.</span></h3>")
	} else {
		fmt.Fprintf(w, "<h3>Failure, something went wrong</h3>")
	}
	fmt.Fprintf(w, "</body></html>")
}
