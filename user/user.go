package user

import (
	apiModels "github.com/stevezaluk/mtgjson-models/api"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	userModel "github.com/stevezaluk/mtgjson-models/user"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"net/http"
)

/*
UserAPI - A representation of the user namespace for the MTGJSON API
*/
type UserAPI struct {
	// baseUrl - The baseUrl with its associating endpoint attached to it, used for making HTTP requests
	baseUrl string

	// client - A pointer to the client.HTTPClient structure that is used for HTTP requests
	client *client.HTTPClient
}

/*
New Create a new instance of the UserAPI struct
*/
func New(baseUrl string, client *client.HTTPClient) *UserAPI {
	return &UserAPI{
		baseUrl: baseUrl,
		client:  client,
	}
}

/*
GetUser Fetch a user based on their email address. Returns ErrNoUser if the user cannot be found
and ErrInvalidEmail if an empty string or invalid email address is passed in the parameter
*/
func (api *UserAPI) GetUser(email string) (*userModel.User, error) {
	request := api.client.BuildRequest(&userModel.User{}).SetQueryParam("email", email)

	resp, err := request.Get(api.baseUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, sdkErrors.ErrTokenInvalid
	}

	if resp.StatusCode() == http.StatusForbidden {
		return nil, sdkErrors.ErrInvalidPermissions
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, sdkErrors.ErrNoUser
	}

	if resp.StatusCode() == http.StatusBadRequest {
		return nil, sdkErrors.ErrInvalidEmail
	}

	return resp.Result().(*userModel.User), nil
}

/*
DeactivateUser Completely removes the requested user account, both from Auth0 and from MongoDB
*/
func (api *UserAPI) DeactivateUser(email string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).SetQueryParam("email", email)

	resp, err := request.Delete(api.baseUrl)
	if err != nil {
		return nil, err
	}

	errResult := resp.Error().(*apiModels.APIResponse)

	if resp.StatusCode() == http.StatusUnauthorized {
		return errResult, sdkErrors.ErrTokenInvalid
	}

	if resp.StatusCode() == http.StatusForbidden {
		return errResult, sdkErrors.ErrInvalidPermissions
	}

	if resp.StatusCode() == http.StatusNotFound {
		return errResult, sdkErrors.ErrNoUser
	}

	if resp.StatusCode() == http.StatusBadRequest {
		return errResult, sdkErrors.ErrInvalidEmail
	}

	return resp.Result().(*apiModels.APIResponse), nil
}
