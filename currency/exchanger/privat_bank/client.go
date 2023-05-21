package privat_bank

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/glbter/currency-ex/currency/exchanger"
)

const (
	pbFormat = "02.01.2006" // DD.MM.YYYY
)

// https://api.privatbank.ua/p24api/exchange_rates?date=01.12.2014
type PrivatBankCurrencyRater struct {
	url    string
	client *http.Client
}

var _ exchanger.CurrencyRater = &PrivatBankCurrencyRater{}

func NewClient(c *http.Client, url string) PrivatBankCurrencyRater {
	return PrivatBankCurrencyRater{
		url:    url,
		client: c,
	}
}

func (cr PrivatBankCurrencyRater) FindRate(c exchanger.Currency, date time.Time) (exchanger.CurrencyRate, error) {
	resp, err := http.Get(cr.url + "?date=" + date.Format(pbFormat))
	if err != nil {
		return exchanger.CurrencyRate{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return exchanger.CurrencyRate{}, fmt.Errorf("responded with %v http code", resp.StatusCode)
	}

	var r pbResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return exchanger.CurrencyRate{}, err
	}

	rate, ok := r.FindRate(c)
	if !ok {
		return exchanger.CurrencyRate{}, errors.New("not found")
	}

	return exchanger.CurrencyRate{
		Base:     exchanger.UAH,
		Rated:    c,
		Sale:     rate.GetSale(),
		Purchase: rate.GetBuy(),
		Date:     date,
	}, nil
}

type pbResp struct {
	Rate []pbRespRate `json:"exchangeRate"`
}

func (r pbResp) FindRate(c exchanger.Currency) (pbRespRate, bool) {
	for _, rate := range r.Rate {
		if rate.Currency == string(c) {
			return rate, true
		}
	}

	return pbRespRate{}, false
}

type pbRespRate struct {
	Base     string   `json:"baseCurrency"`
	Currency string   `json:"currency"`
	Sale     *float64 `json:"saleRate,omitempty"`
	Buy      *float64 `json:"purchaseRate,omitempty"`
	SaleNB   float64  `json:"saleRateNB"`
	BuyNB    float64  `json:"purchaseRateNB"`
}

func (rr pbRespRate) GetSale() float64 {
	if rr.Sale != nil {
		return *rr.Sale
	}

	return rr.SaleNB
}

func (rr pbRespRate) GetBuy() float64 {
	if rr.Buy != nil {
		return *rr.Buy
	}

	return rr.BuyNB
}
