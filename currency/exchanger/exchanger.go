package exchanger

import (
	"fmt"
	"time"
)

type CurrencyRater interface {
	FindRate(c Currency, date time.Time) (CurrencyRate, error)
}

//mockgen -package mock -destination currency/exchanger/mock/mock.go  github.com/glbter/currency-ex/currency/exchanger CurrencyRater

type CurrencySeriesRater interface {
	FindRates(c Currency, start time.Time, end time.Time) ([]CurrencyRate, error)
}

type CurrencyRate struct {
	Base     Currency
	Rated    Currency
	Sale     float64
	Purchase float64
	Date     time.Time
}

type Currency string

const (
	UAH Currency = "UAH"
	EUR Currency = "EUR"
	USD Currency = "USD"
)

func ToCurrency(c string) (Currency, error) {
	switch c {
	case string(UAH):
		return UAH, nil
	case string(USD):
		return USD, nil
	case string(EUR):
		return EUR, nil
	default:
		return "", fmt.Errorf("unknown currency %v", c)
	}
}
