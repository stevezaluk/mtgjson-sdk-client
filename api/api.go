package api

import (
	"github.com/spf13/viper"
	"github.com/stevezaluk/mtgjson-sdk-client/auth"
	"github.com/stevezaluk/mtgjson-sdk-client/card"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"github.com/stevezaluk/mtgjson-sdk-client/deck"
)

/*
MtgjsonApi A representation of the MTGJSON API and all of its routes
*/
type MtgjsonApi struct {
	Client *client.HTTPClient
	Card   *card.CardApi
	Deck   *deck.DeckApi
	Auth   *auth.AuthApi
}

/*
New Initialize a new MTGJSON API object
*/
func New() *MtgjsonApi {
	httpClient := client.NewHttpClient()

	baseUrl := viper.GetString("api.base_url")
	return &MtgjsonApi{
		Client: httpClient,
		Card:   card.New(baseUrl+"/card", httpClient),
		Deck:   deck.New(baseUrl+"/deck", httpClient),
		Auth:   auth.New(httpClient),
	}
}
