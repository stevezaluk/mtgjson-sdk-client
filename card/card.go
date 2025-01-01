package card

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/stevezaluk/mtgjson-models/card"
	"github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk-client/context"
)

func GetCard(uuid string) (card.CardSet, error) {
	var result card.CardSet

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

func NewCard(card card.CardSet) error {
	var uri = context.GetUri("/card")

	cardBytes, err := json.Marshal(&card)
	if err != nil {
		return err
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(cardBytes))
	if err != nil {
		return err
	}

	if resp.StatusCode == 409 {
		return errors.ErrCardAlreadyExist
	}

	if resp.StatusCode == 400 {
		return errors.ErrCardMissingId
	}

	if resp.StatusCode == 500 {
		return errors.ErrCardUpdateFailed
	}

	return nil
}

func DeleteCard(uuid string) error {
	var uri = context.GetUri("/card") + "?cardId=" + uuid

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return errors.ErrNoCard
	}

	if resp.StatusCode == 500 {
		return errors.ErrCardDeleteFailed
	}

	return nil
}

func IndexCards(limit int) ([]card.CardSet, error) {
	var result []card.CardSet

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
