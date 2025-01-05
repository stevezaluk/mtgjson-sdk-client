package set

import (
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
