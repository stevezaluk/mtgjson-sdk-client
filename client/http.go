package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

/*
HTTPClient Simple abstraction of an HTTP Client that can be passed in between the namespaces
of each API
*/
type HTTPClient struct {
	Client *resty.Client
}

/*
NewHttpClient Constructor function for building a new HTTP Client. This should get called once
and then passed between each namespace of the API
*/
func NewHttpClient() *HTTPClient {
	return &HTTPClient{
		Client: resty.New(),
	}
}

/*
BuildRequest Builds a new resty request automatically, filling in the headers and the authentication token
*/
func (client *HTTPClient) BuildRequest(result interface{}) *resty.Request {
	request := client.Client.R().
		EnableTrace().
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "MTGJSON-SDK-Client v1.0.0").
		SetResult(result).
		SetAuthToken(viper.GetString("api.token")) // request will fail if token is not valid

	return request
}
