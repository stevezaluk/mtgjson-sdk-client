package user

import "github.com/stevezaluk/mtgjson-sdk-client/client"

/*
UserApi A representation of the user namespace for the MTGJSON API
*/
type UserApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the UserApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *UserApi {
	return &UserApi{
		BaseUrl: baseUrl,
		client:  client,
	}
}
