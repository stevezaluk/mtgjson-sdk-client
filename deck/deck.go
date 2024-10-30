package deck

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/stevezaluk/mtgjson-models/deck"
	"github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk-client/context"
)

func GetDeck(code string) (deck.Deck, error) {
	var result deck.Deck

	var uri = context.GetUri("/deck") + "?deckCode=" + code

	resp, err := http.Get(uri)

	if resp.StatusCode == 404 {
		return result, errors.ErrNoDeck
	}

	if err != nil {
		return result, err
	}

	body, _ := io.ReadAll(resp.Body)

	if _err := json.Unmarshal(body, &result); err != nil {
		return result, _err
	}

	return result, nil
}
