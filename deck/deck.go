package deck

import (
	deckModel "github.com/stevezaluk/mtgjson-models/deck"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"net/http"
)

/*
DeckApi A representation of the deck namespace for the MTGJSON API
*/
type DeckApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the DeckApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *DeckApi {
	// add check to validate baseUrl here

	return &DeckApi{
		BaseUrl: baseUrl,
		client:  client,
	}
}

/*
GetDeck Fetch a deck from the MongoDB database using the code passed in the parameter. Owner
is the email address of the user that you want to assign to the deck. If the string is empty
then it does not filter by user. Returns ErrNoDeck if the deck does not exist or cannot be located
*/
func (api *DeckApi) GetDeck(code string, owner string) (*deckModel.Deck, error) {
	request := api.client.BuildRequest(&deckModel.Deck{}).
		SetQueryParams(map[string]string{"deckCode": code, "owner": owner})

	resp, err := request.Get(api.BaseUrl)
	if err != nil {
		return nil, err
	}

	if resp.Error() != nil {
		if resp.StatusCode() == http.StatusUnauthorized {
			return nil, sdkErrors.ErrTokenInvalid
		}

		if resp.StatusCode() == http.StatusForbidden {
			return nil, sdkErrors.ErrInvalidPermissions
		}

		if resp.StatusCode() == http.StatusNotFound {
			return nil, sdkErrors.ErrNoDeck
		}

		if resp.StatusCode() == http.StatusBadRequest { // this should never get returned
			return nil, sdkErrors.ErrDeckMissingId
		}
	}

	return resp.Result().(*deckModel.Deck), nil

}
