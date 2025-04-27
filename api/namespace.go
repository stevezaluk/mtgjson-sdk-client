package api

import (
	"github.com/stevezaluk/mtgjson-sdk-client/client"
)

/*
Namespace - An interface that all API namespaces (Card, Deck, User, etc.) implement
*/
type Namespace interface {
	// Client - A function that returns a pointer to the HTTP client used for the namespace
	Client() *client.HTTPClient

	// BaseURL - A function that returns the base URL used in HTTP requests
	BaseURL() string
}
