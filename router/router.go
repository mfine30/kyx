package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/strava/go.strava"
	"github.com/tedsuo/rata"
)

func NewRouter(authenticator *strava.OAuthAuthenticator) (http.Handler, error) {
	handlers := rata.Handlers{
		"root":  newIndexHandler(authenticator),
		"oauth": newoAuthHandler(authenticator),
	}

	callBackPath, err := authenticator.CallbackPath()
	if err != nil {
		fmt.Print("here")
		return nil, err
	}

	routes := rata.Routes{
		{Name: "root", Method: "GET", Path: "/"},
		{Name: "oauth", Method: "GET", Path: callBackPath},
	}

	router, err := rata.NewRouter(routes, handlers)
	if err != nil {
		fmt.Print("there")
		return nil, err
	}

	return router, nil
}

func newIndexHandler(authenticator *strava.OAuthAuthenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// you should make this a template in your real application
		fmt.Fprintf(w, `<a href="%s">`, authenticator.AuthorizationURL("state1", strava.Permissions.Public, true))
		fmt.Fprint(w, `<img src="http://strava.github.io/api/images/ConnectWithStrava.png" />`)
		fmt.Fprint(w, `</a>`)
	})
}

func newoAuthHandler(authenticator *strava.OAuthAuthenticator) http.Handler {
	return authenticator.HandlerFunc(oAuthSuccess, oAuthFailure)
}

func oAuthSuccess(auth *strava.AuthorizationResponse, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "SUCCESS:\nAt this point you can use this information to create a new user or link the account to one of your existing users\n")
	fmt.Fprintf(w, "State: %s\n\n", auth.State)
	fmt.Fprintf(w, "Access Token: %s\n\n", auth.AccessToken)

	fmt.Fprintf(w, "The Authenticated Athlete (you):\n")
	content, _ := json.MarshalIndent(auth.Athlete, "", " ")
	fmt.Fprint(w, string(content))
}

func oAuthFailure(err error, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authorization Failure:\n")

	// some standard error checking
	if err == strava.OAuthAuthorizationDeniedErr {
		fmt.Fprint(w, "The user clicked the 'Do not Authorize' button on the previous page.\n")
		fmt.Fprint(w, "This is the main error your application should handle.")
	} else if err == strava.OAuthInvalidCredentialsErr {
		fmt.Fprint(w, "You provided an incorrect client_id or client_secret.\nDid you remember to set them at the begininng of this file?")
	} else if err == strava.OAuthInvalidCodeErr {
		fmt.Fprint(w, "The temporary token was not recognized, this shouldn't happen normally")
	} else if err == strava.OAuthServerErr {
		fmt.Fprint(w, "There was some sort of server error, try again to see if the problem continues")
	} else {
		fmt.Fprint(w, err)
	}
}
