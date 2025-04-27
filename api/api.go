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
MtgjsonApi A representation of the MTGJSON API and all of its routes
*/
type MtgjsonApi struct {
	Client *client.HTTPClient
	Card   *card.CardApi
	Deck   *deck.DeckApi
	Set    *set.SetApi
	Auth   *auth.AuthApi
	User   *user.UserApi
}

/*
New - Construct a new MtgjsonAPI structure using a hostname and port. If useSSL is set
to true then the protocol will be switched to HTTPS
*/
func New(hostname string, port int, useSSL bool) *MtgjsonApi {
	httpClient := client.NewHttpClient()

	protocol := "http://"
	if useSSL {
		protocol = "https://"
	}

	baseUrl := protocol + hostname + ":" + strconv.Itoa(port)

	return &MtgjsonApi{
		Client: httpClient,
		Card:   card.New(baseUrl+"/card", httpClient),
		Deck:   deck.New(baseUrl+"/deck", httpClient),
		Set:    set.New(baseUrl+"/set", httpClient),
		Auth:   auth.New(httpClient),
		User:   user.New(baseUrl+"/user", httpClient),
	}
}

/*
FromConfig - Construct a new MtgjsonAPI structure using viper config values
*/
func FromConfig() *MtgjsonApi {
	return New(
		viper.GetString("api.hostname"),
		viper.GetInt("api.port"),
		viper.GetBool("api.use_ssl"),
	)
}
