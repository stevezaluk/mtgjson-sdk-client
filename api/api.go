package api

/*
MtgjsonApi A representation of the MTGJSON API and all of its routes
*/
type MtgjsonApi struct {
	Client *HTTPClient
}

/*
New Initialize a new MTGJSON API object
*/
func New() *MtgjsonApi {
	return &MtgjsonApi{
		Client: NewHttpClient(),
	}
}
