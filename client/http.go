package client

import (
	"github.com/go-resty/resty/v2"
	apiModels "github.com/stevezaluk/mtgjson-models/api"
)

/*
HTTPClient Simple abstraction of an HTTP Client that can be passed in between the namespaces
of each API
*/
type HTTPClient struct {
	Client *resty.Client
}

/*
New Constructor function for building a new HTTP Client. This should get called once
and then passed between each namespace of the API
*/
func New() *HTTPClient {
	return &HTTPClient{
		Client: resty.New(),
	}
}

/*
SetBearerToken - Sets the authentication token for the current session
*/
func (client *HTTPClient) SetBearerToken(token string, request *resty.Request) {
	request.SetAuthToken(token)
}

/*
BuildRequest Builds a new resty request automatically, filling in the headers and the authentication token
*/
func (client *HTTPClient) BuildRequest(result interface{}) *resty.Request {
	request := client.Client.R().
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "MTGJSON-SDK-Client v1.0.0").
		SetResult(result).
		SetError(&apiModels.APIResponse{})

	return request
}
