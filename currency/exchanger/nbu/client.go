package nbu

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/glbter/currency-ex/currency/exchanger"
)

const (
	nbuFormat = "20060201" // YYYYMMDD
)

// https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?date=20100302&json

type NBUCurrencyRater struct {
	url    string
	client *http.Client
}

var _ exchanger.CurrencyRater = &NBUCurrencyRater{}

func NewClient(c *http.Client, url string) NBUCurrencyRater {
	return NBUCurrencyRater{
		url:    url,
		client: c,
	}
}

func (cr NBUCurrencyRater) FindRate(c exchanger.Currency, date time.Time) (exchanger.CurrencyRate, error) {
	resp, err := cr.client.Get(cr.url + "?json&date=" + date.Format(nbuFormat))
	if err != nil {
		return exchanger.CurrencyRate{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return exchanger.CurrencyRate{}, fmt.Errorf("responded with %v http code", resp.StatusCode)
	}

	var r []nbuResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return exchanger.CurrencyRate{}, err
	}

	for _, rate := range r {
		if rate.Currency == string(c) {
			return exchanger.CurrencyRate{
				Base:     exchanger.UAH,
				Rated:    c,
				Sale:     rate.Rate,
				Purchase: rate.Rate,
				Date:     date,
			}, nil
		}
	}

	return exchanger.CurrencyRate{}, errors.New("not found")
}

type nbuResp struct {
	UkrName  string  `json:"txt"`
	Rate     float64 `json:"rate"`
	Currency string  `json:"cc"`
}
