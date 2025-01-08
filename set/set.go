package set

import (
	apiModels "github.com/stevezaluk/mtgjson-models/api"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	setModel "github.com/stevezaluk/mtgjson-models/set"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"net/http"
)

/*
SetApi A representation of the set namespace for the MTGJSON API
*/
type SetApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the SetApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *SetApi {
	return &SetApi{
		BaseUrl: baseUrl,
		client:  client,
	}
}

/*
GetSet Takes a single string representing a set code and returns a set model for the set.
Returns ErrNoSet if the set does not exist, or cannot be located
*/
func (api *SetApi) GetSet(code string, owner string) (*setModel.Set, error) {
	request := api.client.BuildRequest(&setModel.Set{}).
		SetQueryParams(map[string]string{"setCode": code, "owner": owner})

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
			return nil, sdkErrors.ErrNoSet
		}
	}

	return resp.Result().(*setModel.Set), nil
}

/*
IndexSets Returns all sets in the database unmarshalled as card models. The limit parameter
will be passed directly to the database query to limit the number of models returned
*/
func (api *SetApi) IndexSets(limit int) (*[]*setModel.Set, error) {
	request := api.client.BuildRequest(&[]*setModel.Set{})

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

		if resp.StatusCode() == http.StatusBadRequest {
			return nil, sdkErrors.ErrNoSets
		}
	}

	return resp.Result().(*[]*setModel.Set), nil
}

/*
NewSet Insert a new set in the form of a model into the MongoDB database. The set model must have a
valid name and set code, additionally the set cannot already exist under the same set code. Owner is
the email address of the owner you want to assign the deck to. If the string is empty (i.e. == ""), it
will be assigned to the system user
*/
func (api *SetApi) NewSet(set *setModel.Set, owner string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).SetQueryParam("owner", owner).SetBody(set)

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
			return errorResponse, sdkErrors.ErrSetAlreadyExists
		}

		if resp.StatusCode() == http.StatusBadRequest {
			if errorResponse.Err == sdkErrors.ErrInvalidObjectStructure.Error() {
				return errorResponse, sdkErrors.ErrInvalidObjectStructure
			}

			if errorResponse.Err == sdkErrors.ErrSetMissingId.Error() {
				return errorResponse, sdkErrors.ErrSetMissingId
			}

			if errorResponse.Err == sdkErrors.ErrMetaApiMustBeNull.Error() {
				return errorResponse, sdkErrors.ErrMetaApiMustBeNull
			}

			if errorResponse.Err == sdkErrors.ErrInvalidCards.Error() {
				return errorResponse, sdkErrors.ErrInvalidCards
			}
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}

/*
DeleteSet Remove a set from the MongoDB database using the code passed in the parameter.
Returns ErrNoSet if the set does not exist. Returns ErrSetDeleteFailed if the deleted count
does not equal 1
*/
func (api *SetApi) DeleteSet(code string, owner string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).SetQueryParams(map[string]string{"setCode": code, "owner": owner})

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
			return errorResponse, sdkErrors.ErrNoSet
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			return errorResponse, sdkErrors.ErrSetDeleteFailed
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}
