package user

import (
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	userModel "github.com/stevezaluk/mtgjson-models/user"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"net/http"
)

/*
UserApi A representation of the user namespace for the MTGJSON API
*/
type UserApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the UserApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *UserApi {
	return &UserApi{
		BaseUrl: baseUrl,
		client:  client,
	}
}

/*
GetUser Fetch a user based on their email addres. Returns ErrNoUser if the user cannot be found
and ErrInvalidEmail if an empty string or invalid email address is passed in the paramter
*/
func (api *UserApi) GetUser(email string) (*userModel.User, error) {
	request := api.client.BuildRequest(&userModel.User{}).SetQueryParam("email", email)

	resp, err := request.Get(api.BaseUrl)
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
