package exchanger

import (
	"context"
	"fmt"
	"time"
)

type AllCurrencyRater interface {
	FindRates(ctx context.Context, params ConvertCurrencyParams) ([]CurrencyRate, error)
}

type ConvertCurrencyParams struct {
	ConvertFrom Currency
	ConvertTo   Currency
	Start       time.Time
	End         time.Time
}

type CurrencyRater interface {
	FindRate(c Currency, date time.Time) (CurrencyRate, error)
}

type CurrencySeriesRater interface {
	FindRates(ctx context.Context, c Currency, start time.Time, end time.Time) ([]CurrencyRate, error)
}

type FindExchangeRateParams struct {
	Currency Currency
	Start    *time.Time
	End      *time.Time
}

type CurrencyRate struct {
	Base     Currency
	Rated    Currency
	Sale     float64
	Purchase float64
	Date     time.Time
	Source   ExchangeSource
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

type ExchangeSource string

const (
	ExchangeSourcePrivat ExchangeSource = "Privat"
	ExchangeSourceNBU    ExchangeSource = "NBU"
)
