package deck

import (
	"github.com/stevezaluk/mtgjson-sdk-client/client"
)

/*
DeckApi A representation of the deck namespace for the MTGJSON API
*/
type DeckApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the DeckApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *DeckApi {
	// add check to validate baseUrl here

	return &DeckApi{
		BaseUrl: baseUrl + "/deck",
		client:  client,
	}
}
