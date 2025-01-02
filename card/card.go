package card

import (
	"github.com/stevezaluk/mtgjson-sdk-client/client"
)

/*
CardApi A representation of the card namespace for the MTGJSON API
*/
type CardApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the CardApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *CardApi {
	// add error check for invalid url here

	return &CardApi{
		BaseUrl: baseUrl,
		client:  client,
	}
}
