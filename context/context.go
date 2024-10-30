package context

import (
	"context"

	"github.com/stevezaluk/mtgjson-sdk-client/config"
)

var ClientContext = context.Background()

const (
	CARD_ENDPOINT    = "/card"
	DECK_ENDPOINT    = "/deck"
	SET_ENDPOINT     = "/set"
	HEALTH_ENDPOINT  = "/health"
	METRICS_ENDPOINT = "/metrics"
)

func InitConfig(config config.Config) {
	ctx := context.WithValue(ClientContext, "config", config)
	ClientContext = ctx
}

func InitUri(config config.Config) {
	ctx := context.WithValue(ClientContext, "uri", config.BuildUri())
	ClientContext = ctx
}

func GetUri(endpoint string) string {
	uri := ClientContext.Value("uri")
	if uri == nil {
		config := ClientContext.Value("config").(config.Config)
		InitUri(config)
	}

	ret := uri.(string) + endpoint

	return ret
}
