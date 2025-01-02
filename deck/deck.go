package deck

import (
	"bytes"
	"encoding/json"
	"github.com/stevezaluk/mtgjson-sdk-client/client"
	"io"
	"net/http"
	"strconv"

	"github.com/stevezaluk/mtgjson-models/deck"
	"github.com/stevezaluk/mtgjson-models/errors"
)

/*
DeckApi A representation of the deck namespace for the MTGJSON API
*/
type DeckApi struct {
	BaseUrl string
	client  *client.HTTPClient
}

/*
New Create a new instance of the DeckApi struct
*/
func New(baseUrl string, client *client.HTTPClient) *DeckApi {
	// add check to validate baseUrl here

	return &DeckApi{
		BaseUrl: baseUrl + "/deck",
		client:  client,
	}
}

func (deckApi *DeckApi) GetDeck(code string) (deck.Deck, error) {
	var result deck.Deck // change to pointer here

	var uri = deckApi.BaseUrl + "?deckCode=" + code // update this to use builtin query string args

	resp, err := http.Get(uri)

	if resp.StatusCode == 404 {
		return result, errors.ErrNoDeck
	}

	if err != nil {
		return result, err
	}

	body, _ := io.ReadAll(resp.Body)

	if _err := json.Unmarshal(body, &result); _err != nil {
		return result, _err
	}

	return result, nil
}

func (deckApi *DeckApi) NewDeck(deck deck.Deck) (bool, error) {
	if deck.Name == "" || deck.Code == "" {
		return false, errors.ErrDeckMissingId
	}

	deckBytes, err := json.Marshal(&deck)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(deckApi.BaseUrl, "application/json", bytes.NewBuffer(deckBytes))

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

func (deckApi *DeckApi) CreateDeck(name string, code string, deckType string) (bool, error) {
	var new deck.Deck // change to pointer here

	if name == "" || code == "" {
		return false, errors.ErrDeckMissingId
	}

	new.Name = name
	new.Code = code
	new.Type = deckType

	_, err := deckApi.NewDeck(new)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (deckApi *DeckApi) DeleteDeck(code string) (bool, error) {
	var uri = deckApi.BaseUrl + "?deckCode=" + code

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

func (deckApi *DeckApi) IndexDecks(limit int) ([]deck.Deck, error) {
	if limit == 0 {
		limit = 100
	}

	var uri = deckApi.BaseUrl + "?limit=" + strconv.Itoa(limit)

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 400 { // change this to 404
		return nil, errors.ErrNoDecks
	}

	var results []deck.Deck

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	return results, nil
}
