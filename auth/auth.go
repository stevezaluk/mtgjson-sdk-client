package auth

import "github.com/stevezaluk/mtgjson-sdk-client/client"

/*
AuthApi A representation of the auth namespace for the MTGJSON API
*/
type AuthApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the AuthApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *AuthApi {
	return &AuthApi{
		BaseUrl: baseUrl,
		client:  client,
	}
}
