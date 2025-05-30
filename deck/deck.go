package deck

import (
	apiModels "github.com/stevezaluk/mtgjson-models/api"
	deckModel "github.com/stevezaluk/mtgjson-models/deck"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"net/http"
)

/*
DeckAPI A representation of the deck namespace for the MTGJSON API
*/
type DeckAPI struct {
	// baseUrl - The baseUrl with its associating endpoint attached to it, used for making HTTP requests
	baseUrl string

	// client - A pointer to the client.HTTPClient structure that is used for HTTP requests
	client *client.HTTPClient
}

/*
New Create a new instance of the DeckAPI struct
*/
func New(baseUrl string, client *client.HTTPClient) *DeckAPI {
	// add check to validate baseUrl here

	return &DeckAPI{
		baseUrl: baseUrl + "/deck",
		client:  client,
	}
}

/*
BaseURL - Returns the baseUrl with its associating endpoint attached to it, used for making HTTP requests
*/
func (api *DeckAPI) BaseURL() string {
	return api.baseUrl
}

/*
Client - Returns a pointer to the client.HTTPClient structure that is used for HTTP requests
*/
func (api *DeckAPI) Client() *client.HTTPClient {
	return api.client
}

/*
GetDeck Fetch a deck from the MongoDB database using the code passed in the parameter. Owner
is the email address of the user that you want to assign to the deck. If the string is empty
then it does not filter by user. Returns ErrNoDeck if the deck does not exist or cannot be located
*/
func (api *DeckAPI) GetDeck(code string, owner string) (*deckModel.Deck, error) {
	request := api.client.BuildRequest(&deckModel.Deck{}).
		SetQueryParams(map[string]string{"deckCode": code, "owner": owner})

	resp, err := request.Get(api.baseUrl)
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
NewDeck Insert a new deck in the form of a model into the MongoDB database. The deck model must have a
valid name and deck code, additionally the deck cannot already exist under the same deck code. Owner is
the email address of the owner you want to assign the deck to. If the string is empty, it will be assigned
to the system user
*/
func (api *DeckAPI) NewDeck(deck *deckModel.Deck, owner string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).
		SetQueryParam("owner", owner).
		SetBody(deck)

	resp, err := request.Post(api.baseUrl)
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

		if resp.StatusCode() == http.StatusConflict {
			return errorResponse, sdkErrors.ErrDeckAlreadyExists
		}

		if resp.StatusCode() == http.StatusBadRequest {
			if errorResponse.Err == sdkErrors.ErrMetaApiMustBeNull.Error() {
				return errorResponse, sdkErrors.ErrMetaApiMustBeNull
			}

			if errorResponse.Err == sdkErrors.ErrDeckMissingContentIds.Error() {
				return errorResponse, sdkErrors.ErrDeckMissingContentIds
			}

			if errorResponse.Err == sdkErrors.ErrDeckMissingId.Error() {
				return errorResponse, sdkErrors.ErrDeckMissingId
			}

			if errorResponse.Err == sdkErrors.ErrInvalidCards.Error() {
				return errorResponse, sdkErrors.ErrInvalidCards
			}
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}

/*
DeleteDeck Remove a deck from the MongoDB database using the code passed in the
parameter. Returns ErrNoDeck if the deck does not exist. Returns
ErrDeckDeleteFailed if the deleted count does not equal 1
*/
func (api *DeckAPI) DeleteDeck(code string, owner string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).
		SetQueryParams(map[string]string{"deckCode": code, "owner": owner})

	resp, err := request.Delete(api.baseUrl)
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
			return errorResponse, sdkErrors.ErrDeckMissingId
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			return errorResponse, sdkErrors.ErrDeckDeleteFailed
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}

/*
GetDeckContents Update the 'contents' field of the deck passed in the parameter. This accepts a
pointer and updates this in place to avoid having to copy large amounts of data
*/
func (api *DeckAPI) GetDeckContents(code string, owner string) (*deckModel.DeckContents, error) {
	request := api.client.BuildRequest(&deckModel.DeckContents{}).SetQueryParams(map[string]string{"deckCode": code, "owner": owner})

	resp, err := request.Get(api.baseUrl + "/content")
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
func (api *DeckAPI) AddCards(code string, cards *deckModel.DeckContentIds, owner string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).
		SetQueryParams(map[string]string{"deckCode": code, "owner": owner}).
		SetBody(cards)

	resp, err := request.Post(api.baseUrl + "/content")
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

			if errorResponse.Err == sdkErrors.ErrDeckNoCards.Error() {
				return errorResponse, sdkErrors.ErrDeckNoCards
			}
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			return errorResponse, sdkErrors.ErrDeckUpdateFailed
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}

func (api *DeckAPI) RemoveCards(code string, cards *deckModel.DeckContentIds, owner string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).
		SetQueryParams(map[string]string{"deckCode": code, "owner": owner}).
		SetBody(cards)

	resp, err := request.Delete(api.baseUrl + "/content")
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

		if resp.StatusCode() == http.StatusBadRequest {
			if errorResponse.Err == sdkErrors.ErrDeckMissingId.Error() {
				return errorResponse, sdkErrors.ErrDeckMissingId
			}

			if errorResponse.Err == sdkErrors.ErrInvalidCards.Error() {
				return errorResponse, sdkErrors.ErrInvalidCards
			}

			if errorResponse.Err == sdkErrors.ErrDeckNoCards.Error() {
				return errorResponse, sdkErrors.ErrDeckNoCards
			}
		}

		if resp.StatusCode() == http.StatusNotFound {
			return errorResponse, sdkErrors.ErrNoDeck
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			return errorResponse, sdkErrors.ErrDeckUpdateFailed
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}
