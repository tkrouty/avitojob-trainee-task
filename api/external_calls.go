package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func getExchangeRatebyHTTP(currency string) (float64, error) {
	baseURL, _ := url.Parse("https://api.exchangeratesapi.io/latest")
	params := url.Values{}
	params.Add("base", "RUB")
	params.Add("symbols", currency)
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r struct {
		Date  string             `json:"date"`
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}
	rate := r.Rates[currency]

	if rate == 0 {
		return 0, errors.New("Unable to get exchange rate")
	}

	return rate, nil
}
