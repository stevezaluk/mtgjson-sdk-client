package card

import (
	apiModels "github.com/stevezaluk/mtgjson-models/api"
	cardModel "github.com/stevezaluk/mtgjson-models/card"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"net/http"
)

/*
CardAPI A representation of the card namespace for the MTGJSON API
*/
type CardAPI struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the CardAPI struct
*/
func New(baseUrl string, client *client.HTTPClient) *CardAPI {
	// add error check for invalid url here

	return &CardAPI{
		BaseUrl: baseUrl,
		client:  client,
	}
}

/*
GetCard Takes a single string representing an MTGJSONv4 UUID and return a card model
for it
*/
func (api *CardAPI) GetCard(uuid string, owner string) (*cardModel.CardSet, error) {
	request := api.client.BuildRequest(&cardModel.CardSet{}).SetQueryParams(map[string]string{"cardId": uuid, "owner": owner})

	resp, err := request.Get(api.BaseUrl)
	if err != nil {
		return nil, err
	}

	if resp.Error() != nil {
		errorResponse := resp.Error().(*apiModels.APIResponse)
		if resp.StatusCode() == http.StatusUnauthorized {
			return nil, sdkErrors.ErrTokenInvalid
		}

		if resp.StatusCode() == http.StatusForbidden {
			return nil, sdkErrors.ErrInvalidPermissions
		}

		if resp.StatusCode() == http.StatusNotFound {
			return nil, sdkErrors.ErrNoCard
		}

		if resp.StatusCode() == http.StatusBadRequest {
			if errorResponse.Err == sdkErrors.ErrInvalidUUID.Error() {
				return nil, sdkErrors.ErrInvalidUUID
			}
		}
	}

	return resp.Result().(*cardModel.CardSet), nil
}

/*
IndexCards Returns all cards in the database unmarshalled as card models. The limit parameter
will be passed directly to the database query to limit the number of models returned
*/
func (api *CardAPI) IndexCards() (*[]*cardModel.CardSet, error) {
	request := api.client.BuildRequest(&[]*cardModel.CardSet{})

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
			return nil, sdkErrors.ErrNoCards
		}
	}

	return resp.Result().(*[]*cardModel.CardSet), nil
}

/*
NewCard Insert a new card in the form of a model into the MongoDB database. The card model must have a
valid name and MTGJSONv4 ID, additionally, the card cannot already exist under the same ID
*/
func (api *CardAPI) NewCard(card *cardModel.CardSet, owner string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).
		SetBody(card).
		SetQueryParam("owner", owner)

	resp, err := request.Post(api.BaseUrl)
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
			return errorResponse, sdkErrors.ErrCardAlreadyExist
		}

		if resp.StatusCode() == http.StatusBadRequest {
			return errorResponse, sdkErrors.ErrCardMissingId
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}

/*
DeleteCard Remove a card from the MongoDB database. The UUID passed in the parameter must be a valid MTGJSONv4 ID.
ErrNoCard will be returned if no card exists under the passed UUID, and ErrCardDeleteFailed will be returned
if the deleted count does not equal 1
*/
func (api *CardAPI) DeleteCard(uuid string, owner string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).
		SetQueryParams(map[string]string{"cardId": uuid, "owner": owner})

	resp, err := request.Delete(api.BaseUrl)
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
			return errorResponse, sdkErrors.ErrNoCard
		}

		if resp.StatusCode() == http.StatusBadRequest {
			return errorResponse, sdkErrors.ErrCardMissingId
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			return errorResponse, sdkErrors.ErrCardDeleteFailed
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}
