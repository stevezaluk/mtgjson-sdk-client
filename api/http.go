package api

import "github.com/go-resty/resty/v2"

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
