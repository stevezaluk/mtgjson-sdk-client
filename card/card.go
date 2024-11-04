package card

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/stevezaluk/mtgjson-models/card"
	"github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk-client/context"
)

func GetCard(uuid string) (card.Card, error) {
	var result card.Card

	var uri = context.GetUri("/card") + "?cardId=" + uuid

	resp, err := http.Get(uri)

	if resp.StatusCode == 404 {
		return result, errors.ErrNoCard
	}

	if resp.StatusCode == 400 {
		return result, errors.ErrInvalidUUID
	}

	if err != nil {
		return result, errors.ErrNoCard
	}

	body, _ := io.ReadAll(resp.Body)

	if _err := json.Unmarshal(body, &result); _err != nil {
		return result, _err
	}

	return result, nil
}

func IndexCards(limit int) ([]card.Card, error) {
	var result []card.Card

	if limit == 0 {
		limit = 100
	}

	var uri = context.GetUri("/card") + "?limit=" + strconv.Itoa(limit)

	resp, err := http.Get(uri)

	if resp.StatusCode == 400 {
		return result, errors.ErrNoCards
	}

	if err != nil {
		return result, errors.ErrNoCard
	}

	body, _ := io.ReadAll(resp.Body)

	if _err := json.Unmarshal(body, &result); _err != nil {
		return result, _err
	}

	return result, nil
}
