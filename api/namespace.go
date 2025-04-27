package api

import "net/http"

/*
Namespace - An interface that all API namespaces (Card, Deck, User, etc.) implement
*/
type Namespace interface {
	Client() *http.Client
	BaseURL() string
}
