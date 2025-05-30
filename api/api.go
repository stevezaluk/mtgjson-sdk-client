package api

import (
	"github.com/spf13/viper"
	"github.com/stevezaluk/mtgjson-sdk-client/auth"
	"github.com/stevezaluk/mtgjson-sdk-client/card"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"github.com/stevezaluk/mtgjson-sdk-client/deck"
	"github.com/stevezaluk/mtgjson-sdk-client/set"
	"github.com/stevezaluk/mtgjson-sdk-client/user"
	"strconv"
)

/*
MtgjsonAPI - A representation of the MTGJSON API and all of its routes
*/
type MtgjsonAPI struct {
	// client - A pointer to the client.HTTPClient structure that is used for HTTP requests
	client *client.HTTPClient

	// Card - A pointer to the CardAPI namespace, used for making HTTP requests to the /card endpoint
	Card *card.CardAPI

	// Deck - A pointer to the DeckAPI namespace, used for making HTTP requests to the /deck endpoint
	Deck *deck.DeckAPI

	// Set - A pointer to the SetAPI namespace, used for making HTTP requests to the /set endpoint
	Set *set.SetAPI

	// Auth - A pointer to the AuthAPI namespace, used for making HTTP requests to the /login and /register endpoints
	Auth *auth.AuthAPI

	// User - A pointer to the UserAPI namespace, used for making HTTP requests to the /user endpoint
	User *user.UserAPI
}

/*
New - Construct a new MtgjsonAPI structure using a hostname and port. If useSSL is set
to true then the protocol will be switched to HTTPS
*/
func New(hostname string, port int, useSSL bool) *MtgjsonAPI {
	httpClient := client.New()

	protocol := "http://"
	if useSSL {
		protocol = "https://"
	}

	baseUrl := protocol + hostname + ":" + strconv.Itoa(port)

	return &MtgjsonAPI{
		client: httpClient,
		Card:   card.New(baseUrl, httpClient),
		Deck:   deck.New(baseUrl, httpClient),
		Set:    set.New(baseUrl, httpClient),
		Auth:   auth.New(baseUrl, httpClient),
		User:   user.New(baseUrl, httpClient),
	}
}

/*
FromConfig - Construct a new MtgjsonAPI structure using viper config values
*/
func FromConfig() *MtgjsonAPI {
	return New(
		viper.GetString("api.hostname"),
		viper.GetInt("api.port"),
		viper.GetBool("api.use_ssl"),
	)
}

/*
Client - Returns a pointer to the client.HTTPClient structure that is used for HTTP requests
*/
func (api *MtgjsonAPI) Client() *client.HTTPClient {
	return api.client
}

/*
SetEmailPasswordAuth - Fetches a token for the client.HTTPClient to use in authenticated requests
*/
func (api *MtgjsonAPI) SetEmailPasswordAuth(email string, password string) error {
	tokenSet, err := api.Auth.Login(email, password)
	if err != nil {
		return err
	}

	api.Client().SetBearerToken(tokenSet)

	return nil
}
