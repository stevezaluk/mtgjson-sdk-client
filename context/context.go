package context

import (
	"context"

	"github.com/stevezaluk/mtgjson-sdk-client/config"
)

var ClientContext = context.Background()

func InitConfig(config config.Config) {
	ctx := context.WithValue(ClientContext, "config", config)
	ClientContext = ctx
}

func InitUri(config config.Config) {
	ctx := context.WithValue(ClientContext, "uri", config.BuildUri())
	ClientContext = ctx
}

func GetUri() string {
	uri := ClientContext.Value("uri")
	if uri == nil {
		config := ClientContext.Value("config").(config.Config)
		InitUri(config)
	}

	return uri.(string)
}
