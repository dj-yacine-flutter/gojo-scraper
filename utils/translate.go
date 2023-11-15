package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/dj-yacine-flutter/gojo-scraper/models"
)

func Translate(client *http.Client, text, sourceLanguage, targetLanguage string) (*models.LibreTranslate, error) {

	URL := "http://localhost:5000/translate"

	overviewText := strings.ToLower(text)
	postData := url.Values{}
	postData.Set("q", overviewText)
	postData.Set("source", sourceLanguage)
	postData.Set("target", targetLanguage)
	postData.Set("format", "text")
	postData.Set("api_key", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")

	dataEncoded := postData.Encode()

	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(dataEncoded))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	lt := models.LibreTranslate{}
	err = json.NewDecoder(resp.Body).Decode(&lt)
	if err != nil {
		return nil, err
	}

	return &lt, nil
}
