package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CurrencyData struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func FetchRates(apiKey string) (CurrencyData, error) {
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s&prettyprint=false", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return CurrencyData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return CurrencyData{}, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var data CurrencyData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return CurrencyData{}, err
	}

	return data, nil
}
