package api

import "net/http"

/*
Namespace - An interface that all API namespaces (Card, Deck, User, etc.) implement
*/
type Namespace interface {
	// Client - A function that returns a pointer to the HTTP client used for the namespace
	Client() *http.Client

	// BaseURL - A function that returns the base URL used in HTTP requests
	BaseURL() string
}
