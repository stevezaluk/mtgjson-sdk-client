package auth

import (
	"github.com/auth0/go-auth0/authentication/oauth"
	"github.com/spf13/viper"
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
