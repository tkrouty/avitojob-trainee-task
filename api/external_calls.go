package api

import (
    "net/url"
    "net/http"
    "encoding/json"
    "errors"
)

func getExchangeRate(currency string) (float64, error) {
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
