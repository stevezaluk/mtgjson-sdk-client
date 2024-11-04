package health

import (
	"net/http"

	"github.com/stevezaluk/mtgjson-sdk-client/context"
)

func GetHealth() (bool, error) {
	var uri = context.GetUri("/health")

	resp, err := http.Get(uri)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		return true, nil
	}

	return false, nil
}
