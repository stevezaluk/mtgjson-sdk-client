package api

import (
	"github.com/go-resty/resty/v2"
)

/*
MtgjsonApi A representation of the MTGJSON API and all of its routes
*/
type MtgjsonApi struct {
	client *resty.Client
}

/*
New Initialize a new MTGJSON API object
*/
func New() *MtgjsonApi {
	client := resty.New()

	return &MtgjsonApi{
		client: client,
	}
}
