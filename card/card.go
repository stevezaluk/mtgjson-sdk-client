package card

import (
	apiModels "github.com/stevezaluk/mtgjson-models/api"
	cardModel "github.com/stevezaluk/mtgjson-models/card"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"net/http"
)

/*
CardApi A representation of the card namespace for the MTGJSON API
*/
type CardApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the CardApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *CardApi {
	// add error check for invalid url here

	return &CardApi{
		BaseUrl: baseUrl,
		client:  client,
	}
}

/*
GetCard Takes a single string representing an MTGJSONv4 UUID and return a card model
for it
*/
func (api *CardApi) GetCard(uuid string, owner string) (*cardModel.CardSet, error) {
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
