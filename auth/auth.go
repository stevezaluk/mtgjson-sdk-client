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
AuthApi A representation of the auth namespace for the MTGJSON API
*/
type AuthApi struct {
	client *client.HTTPClient
}

/*
New Create a new instance of the AuthApi struct
*/
func New(client *client.HTTPClient) *AuthApi {
	return &AuthApi{
		client: client,
	}
}

/*
Login Exchange user credentials for an oauth.TokenSet
*/
func (api *AuthApi) Login(email string, password string) (*oauth.TokenSet, error) {
	request := api.client.BuildRequest(&oauth.TokenSet{}).
		SetBody(map[string]string{"email": email, "password": password}) // create protobuf model for this

	resp, err := request.Post(viper.GetString("api.base_url") + "/login")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, sdkErrors.ErrNoUser
	}

	if resp.StatusCode() == http.StatusInternalServerError {
		return nil, sdkErrors.ErrTokenInvalid
	}

	return resp.Result().(*oauth.TokenSet), nil
}

/*
SetAuthToken Make a login request for the user and set the auth token for this session
*/
func (api *AuthApi) SetAuthToken(email string, password string) error {
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
func (api *AuthApi) RegisterUser(email string, username string, password string) (*apiModels.APIResponse, error) {
	if email == "" || username == "" || password == "" {
		return nil, sdkErrors.ErrUserMissingId
	}

	request := api.client.BuildRequest(&apiModels.APIResponse{}).
		SetBody(apiModels.RegisterRequest{
			Username: username,
			Email:    email,
			Password: password,
		})

	resp, err := request.Post(viper.GetString("api.base_url") + "/register")
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
func (api *AuthApi) ResetUserPassword(email string) (*apiModels.APIResponse, error) {
	request := api.client.BuildRequest(&apiModels.APIResponse{}).SetQueryParam("email", email)

	resp, err := request.Get(viper.GetString("api.base_url") + "/reset")
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
