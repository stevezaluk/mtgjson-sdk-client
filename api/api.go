package api

import "github.com/stevezaluk/mtgjson-sdk-client/client"

/*
MtgjsonApi A representation of the MTGJSON API and all of its routes
*/
type MtgjsonApi struct {
	Client *client.HTTPClient
}

/*
New Initialize a new MTGJSON API object
*/
func New() *MtgjsonApi {
	return &MtgjsonApi{
		Client: client.NewHttpClient(),
	}
}
