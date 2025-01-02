package card

import (
	"bytes"
	"encoding/json"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"io"
	"net/http"
	"strconv"

	"github.com/stevezaluk/mtgjson-models/card"
	"github.com/stevezaluk/mtgjson-models/errors"
)

/*
CardApi A representation of the card namespace for the MTGJSON API
*/
type CardApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the CardApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *CardApi {
	// add error check for invalid url here

	return &CardApi{
		BaseUrl: baseUrl + "/card",
		client:  client,
	}
}

func (cardApi *CardApi) GetCard(uuid string) (card.CardSet, error) {
	var result card.CardSet // change return type to pointers here

	var uri = cardApi.BaseUrl + "?cardId=" + uuid

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

func (cardApi *CardApi) NewCard(card card.CardSet) error {
	cardBytes, err := json.Marshal(&card)
	if err != nil {
		return err
	}

	resp, err := http.Post(cardApi.BaseUrl, "application/json", bytes.NewBuffer(cardBytes))
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

func (cardApi *CardApi) DeleteCard(uuid string) error {
	req, err := http.NewRequest("DELETE", cardApi.BaseUrl, nil)
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

func (cardApi *CardApi) IndexCards(limit int) ([]card.CardSet, error) {
	var result []card.CardSet

	if limit == 0 {
		limit = 100
	}

	var uri = cardApi.BaseUrl + "?limit=" + strconv.Itoa(limit)

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
