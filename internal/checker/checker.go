package checker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const serviceURL = "https://speller.yandex.net/services/spellservice.json/checkTexts"

type YandexSpellResponse struct {
	Code        int      `json:"code"`
	Pos         int      `json:"pos"`
	Row         int      `json:"row"`
	Col         int      `json:"col"`
	Len         int      `json:"len"`
	Word        string   `json:"word"`
	Suggestions []string `json:"s"`
}

type YandexSpellChecker struct {
	URL string
}

type YandexSpellCheckerResults struct {
}

func New() *YandexSpellChecker {
	return &YandexSpellChecker{
		URL: serviceURL,
	}
}

func (c *YandexSpellChecker) FixSpellsInBatchMode(rows []string) ([]string, error) {
	processed := strings.Join(rows, "&text=")
	formatted := strings.Join(strings.Split(processed, " "), "+")

	response, err := http.Get(c.URL + "?text=" + formatted)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var targets [][]YandexSpellResponse
	if err = json.Unmarshal(body, &targets); err != nil {
		return nil, err
	}

	for i, text := range targets {
		if len(text) > 0 {
			for _, word := range text {
				rows[i] = strings.Replace(rows[i], word.Word, word.Suggestions[0], -1)
			}
		}

	}

	return rows, nil
}
