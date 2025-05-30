package client

import (
	"github.com/auth0/go-auth0/authentication/oauth"
	"github.com/go-resty/resty/v2"
	apiModels "github.com/stevezaluk/mtgjson-models/api"
)

/*
HTTPClient Simple abstraction of an HTTP Client that can be passed in between the namespaces
of each API
*/
type HTTPClient struct {
	// client - The resty.Client structure that is shared across the API namespaces
	client *resty.Client

	// token - The JWT Token Set used for authentication
	token *oauth.TokenSet
}

/*
New Constructor function for building a new HTTP Client. This should get called once
and then passed between each namespace of the API
*/
func New() *HTTPClient {
	return &HTTPClient{
		client: resty.New(),
	}
}

/*
Client - Returns a pointer to the resty.Client structure that is shared across API namespaces
*/
func (client *HTTPClient) Client() *resty.Client {
	return client.client
}

/*
SetBearerToken - Sets the authentication token for the current session
*/
func (client *HTTPClient) SetBearerToken(token *oauth.TokenSet) {
	if token == nil {
		return
	}

	client.token = token

}

/*
BuildRequest Builds a new resty request automatically, filling in the headers and the authentication token
*/
func (client *HTTPClient) BuildRequest(result interface{}) *resty.Request {
	request := client.Client().R().
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "MTGJSON-SDK-Client v1.0.0").
		SetResult(result).
		SetError(&apiModels.APIResponse{})

	if client.token != nil {
		request.SetAuthToken(client.token.AccessToken)
	}

	return request
}
