package api

import (
	"github.com/spf13/viper"
	"github.com/stevezaluk/mtgjson-sdk-client/card"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"github.com/stevezaluk/mtgjson-sdk-client/deck"
)

/*
MtgjsonApi A representation of the MTGJSON API and all of its routes
*/
type MtgjsonApi struct {
	Client *client.HTTPClient
	Deck   *deck.DeckApi
	Card   *card.CardApi
}

/*
New Initialize a new MTGJSON API object
*/
func New() *MtgjsonApi {
	httpClient := client.NewHttpClient()

	baseUrl := viper.GetString("api.base_url")
	return &MtgjsonApi{
		Client: httpClient,
		Deck:   deck.New(baseUrl+"/deck", httpClient),
		Card:   card.New(baseUrl+"/card", httpClient),
	}
}
