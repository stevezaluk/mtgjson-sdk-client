package deck

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

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

	if _err := json.Unmarshal(body, &result); _err != nil {
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

func AddCards(code string, cards []string, board string) ([]string, []string, error) {
	var updates deck.DeckUpdate

	var uri = context.GetUri("/deck/content") + "?deckCode=" + code

	if board == deck.MAINBOARD {
		updates.MainBoard = append(updates.MainBoard, cards...)
	} else if board == deck.SIDEBOARD {
		updates.SideBoard = append(updates.SideBoard, cards...)
	} else if board == deck.COMMANDER {
		updates.Commander = append(updates.Commander, cards...)
	} else {
		return nil, nil, errors.ErrBoardNotExist
	}

	updateBytes, err := json.Marshal(&updates)
	if err != nil {
		return nil, nil, err
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(updateBytes))
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode == 404 {
		return nil, nil, errors.ErrNoDeck
	}

	if resp.StatusCode == 500 {
		return nil, nil, errors.ErrDeckUpdateFailed
	}

	if resp.StatusCode == 400 {

		type InvalidCards struct {
			Invalid []string
			NoExist []string
		}

		var invalidCards InvalidCards
		body, _ := io.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &invalidCards); err != nil {
			return nil, nil, err
		}

		return invalidCards.Invalid, invalidCards.NoExist, errors.ErrDeckUpdateFailed
	}

	return nil, nil, nil
}

func DeleteCards(code string, cards []string, board string) error {
	var updates deck.DeckUpdate

	var uri = context.GetUri("/deck/content") + "?deckCode=" + code

	if board == deck.MAINBOARD {
		updates.MainBoard = append(updates.MainBoard, cards...)
	} else if board == deck.SIDEBOARD {
		updates.SideBoard = append(updates.SideBoard, cards...)
	} else if board == deck.COMMANDER {
		updates.Commander = append(updates.Commander, cards...)
	}

	deleteBytes, err := json.Marshal(&updates)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", uri, bytes.NewBuffer(deleteBytes))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		return errors.ErrNoDeck
	}

	if resp.StatusCode == 500 {
		return errors.ErrDeckDeleteFailed
	}

	return nil
}

func IndexDecks(limit int) ([]deck.Deck, error) {
	if limit == 0 {
		limit = 100
	}

	var uri = context.GetUri("/deck") + "?limit=" + strconv.Itoa(limit)

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
