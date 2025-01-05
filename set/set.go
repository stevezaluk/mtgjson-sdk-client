package set

import "github.com/stevezaluk/mtgjson-sdk-client/client"

/*
SetApi A representation of the set namespace for the MTGJSON API
*/
type SetApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the SetApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *SetApi {
	return &SetApi{
		BaseUrl: baseUrl,
		client:  client,
	}
}
