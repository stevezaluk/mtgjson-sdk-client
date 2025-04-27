package auth

import (
	"errors"
	"github.com/auth0/go-auth0/authentication/oauth"
	"github.com/spf13/viper"
	apiModels "github.com/stevezaluk/mtgjson-models/api"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"net/http"
)

/*
AuthAPI A representation of the auth namespace for the MTGJSON API
*/
type AuthAPI struct {
	// baseUrl - The baseUrl with its associating endpoint attached to it, used for making HTTP requests
	baseUrl string

	// client - A pointer to the client.HTTPClient structure that is used for HTTP requests
	client *client.HTTPClient
}

/*
New Create a new instance of the AuthAPI struct
*/
func New(baseUrl string, client *client.HTTPClient) *AuthAPI {
	return &AuthAPI{
		baseUrl: baseUrl,
		client:  client,
	}
}

/*
BaseURL - Returns the baseUrl with its associating endpoint attached to it, used for making HTTP requests
*/
func (api *AuthAPI) BaseURL() string {
	return api.baseUrl
}

/*
Client - Returns a pointer to the client.HTTPClient structure that is used for HTTP requests
*/
func (api *AuthAPI) Client() *client.HTTPClient {
	return api.client
}

/*
Login Exchange user credentials for an oauth.TokenSet
*/
func (api *AuthAPI) Login(email string, password string) (*oauth.TokenSet, error) {
	request := api.client.BuildRequest(&oauth.TokenSet{}).
		SetBody(apiModels.LoginRequest{
			Email:    email,
			Password: password,
		})

	resp, err := request.Post(api.baseUrl + "/login")
	if err != nil {
		return nil, err
	}

	if resp.Error() != nil { // were not returning the API response here to avoid having more than two responses
		if resp.StatusCode() == http.StatusNotFound {
			return nil, sdkErrors.ErrNoUser
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			return nil, sdkErrors.ErrTokenInvalid
		}
	}

	return resp.Result().(*oauth.TokenSet), nil
}

/*
SetAuthToken Make a login request for the user and set the auth token for this session
*/
func (api *AuthAPI) SetAuthToken(email string, password string) error {
	tokenSet, err := api.Login(email, password)
	if err != nil {
		return err
	}

	viper.Set("api.token_str", tokenSet.AccessToken)
	viper.Set("api.token", tokenSet)

	return nil
}

/*
RegisterUser Register a new user with Auth0 and store there user model within the MongoDB database
*/
func (api *AuthAPI) RegisterUser(email string, username string, password string) (*apiModels.APIResponse, error) {
	if email == "" || username == "" || password == "" {
		return nil, sdkErrors.ErrUserMissingId
	}

	request := api.client.BuildRequest(&apiModels.APIResponse{}).
		SetBody(apiModels.RegisterRequest{
			Username: username,
			Email:    email,
			Password: password,
		})

	resp, err := request.Post(api.baseUrl + "/register")
	if err != nil {
		return nil, err
	}

	if resp.Error() != nil {
		errorResult := resp.Error().(*apiModels.APIResponse)
		if resp.StatusCode() == http.StatusConflict { // this needs to be added to the API
			return errorResult, sdkErrors.ErrUserAlreadyExist
		}

		if resp.StatusCode() == http.StatusBadRequest {
			if errorResult.Err == sdkErrors.ErrInvalidPasswordLength.Error() {
				return errorResult, sdkErrors.ErrInvalidPasswordLength
			}

			if errorResult.Err == sdkErrors.ErrInvalidEmail.Error() {
				return errorResult, sdkErrors.ErrInvalidEmail
			}
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			return errorResult, sdkErrors.ErrFailedToRegisterUser
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}

/*
ResetUserPassword Send a reset password email to a specified user account
*/
func (api *AuthAPI) ResetUserPassword(email string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).SetQueryParam("email", email)

	resp, err := request.Get(api.baseUrl + "/reset")
	if err != nil {
		return nil, err
	}

	if resp.Error() != nil {
		errResult := resp.Error().(*apiModels.APIResponse)

		if resp.StatusCode() == http.StatusUnauthorized {
			return errResult, sdkErrors.ErrTokenInvalid
		}

		if resp.StatusCode() == http.StatusForbidden {
			return errResult, sdkErrors.ErrInvalidPermissions
		}

		if resp.StatusCode() == http.StatusBadRequest {
			return errResult, sdkErrors.ErrInvalidEmail
		}

		if resp.StatusCode() == http.StatusNotFound {
			return errResult, sdkErrors.ErrNoUser
		}

		if resp.StatusCode() == http.StatusInternalServerError {
			return errResult, errors.New("user: Failed to reset user password") // this needs to be added as a named error
		}
	}

	return resp.Result().(*apiModels.APIResponse), nil
}
