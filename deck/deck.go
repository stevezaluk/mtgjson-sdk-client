package deck

import (
	"bytes"
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

func NewDeck(deck deck.Deck) (bool, error) {
	if deck.Name == "" || deck.Code == "" {
		return false, errors.ErrDeckMissingId
	}

	var uri = context.GetUri("/deck")

	deckBytes, err := json.Marshal(&deck)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(deckBytes))

	if resp.StatusCode == 500 {
		return false, errors.ErrDeckUpdateFailed
	}

	if resp.StatusCode == 400 {
		// return ErrDeckInvalid
	}

	if resp.StatusCode == 409 {
		return false, errors.ErrDeckAlreadyExists
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateDeck(name string, code string, deckType string) (bool, error) {
	var new deck.Deck

	if name == "" || code == "" {
		return false, errors.ErrDeckMissingId
	}

	new.Name = name
	new.Code = code
	new.Type = deckType

	_, err := NewDeck(new)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteDeck(code string) (bool, error) {
	var uri = context.GetUri("/deck") + "?deckCode=" + code

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return false, nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, nil
	}

	if resp.StatusCode == 404 {
		return false, errors.ErrNoDeck
	}

	if resp.StatusCode == 500 {
		return false, errors.ErrDeckDeleteFailed
	}

	return true, nil
}
