package deck

import (
	apiModels "github.com/stevezaluk/mtgjson-models/api"
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

/*
GetDeckContents Update the 'contents' field of the deck passed in the parameter. This accepts a
pointer and updates this in place to avoid having to copy large amounts of data
*/
func (api *DeckApi) GetDeckContents(code string, owner string) (*deckModel.DeckContents, error) {
	request := api.client.BuildRequest(&deckModel.DeckContents{}).SetQueryParams(map[string]string{"deckCode": code, "owner": owner})

	resp, err := request.Get(api.BaseUrl + "/content")
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

		if resp.StatusCode() == http.StatusBadRequest {
			return nil, sdkErrors.ErrDeckMissingId
		}
	}

	return resp.Result().(*deckModel.DeckContents), nil
}

/*
AddCards Update the content ids in the deck model passed with new cards. This should
probably validate cards in the future
*/
func (api *DeckApi) AddCards(code string, cards *deckModel.DeckContentIds, owner string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).
		SetQueryParams(map[string]string{"deckCode": code, "owner": owner}).
		SetBody(cards)

	resp, err := request.Post(api.BaseUrl + "/content")
	if err != nil {
		return nil, err
	}

	if resp.Error() != nil {
		errorResponse := resp.Error().(*apiModels.APIResponse)

		if resp.StatusCode() == http.StatusUnauthorized {
			return errorResponse, sdkErrors.ErrTokenInvalid
		}

		if resp.StatusCode() == http.StatusForbidden {
			return errorResponse, sdkErrors.ErrInvalidPermissions
		}

		if resp.StatusCode() == http.StatusNotFound {
			return errorResponse, sdkErrors.ErrNoDeck
		}

		if resp.StatusCode() == http.StatusBadRequest {
			if errorResponse.Err == sdkErrors.ErrInvalidObjectStructure.Error() {
				return errorResponse, sdkErrors.ErrInvalidObjectStructure
			}

			if errorResponse.Err == sdkErrors.ErrDeckMissingContentIds.Error() {
				return errorResponse, sdkErrors.ErrDeckMissingContentIds
			}

			if errorResponse.Err == sdkErrors.ErrInvalidCards.Error() {
				return errorResponse, sdkErrors.ErrInvalidCards
			}
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			return errorResponse, sdkErrors.ErrDeckUpdateFailed
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}
