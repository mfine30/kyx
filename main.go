package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/mfine30/kyx/router"
	"github.com/strava/go.strava"
)

var authenticator *strava.OAuthAuthenticator

func main() {
	host := os.Getenv("KYX_HOST")
	port := os.Getenv("PORT")

	clientId, err := strconv.Atoi(os.Getenv("STRAVA_CLIENT_ID"))
	if err != nil {
		fmt.Printf("Error converting CLIENT_ID to string: %s", err.Error())
	}
	strava.ClientId = clientId
	strava.ClientSecret = os.Getenv("STRAVA_CLIENT_SECRET")

	authenticator = &strava.OAuthAuthenticator{
		CallbackURL:            fmt.Sprintf("http://%s/exchange_token", host),
		RequestClientGenerator: nil,
	}

	router, err := router.NewRouter(authenticator)
	if err != nil {
		panic(err)
	}

	// start the server
	fmt.Printf("Visit http://%s:%s/ to view the demo\n", host, port)
	fmt.Printf("ctrl-c to exit")
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
