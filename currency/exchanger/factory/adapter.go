package factory

import (
	"context"
	"fmt"
	"github.com/glbter/currency-ex/currency/exchanger"
	"time"
)

type AllCurrencyRater struct {
	rater exchanger.CurrencySeriesRater
}

func NewAllCurrencyRater(rater exchanger.CurrencySeriesRater) AllCurrencyRater {
	return AllCurrencyRater{
		rater: rater,
	}
}

func (r AllCurrencyRater) FindRates(ctx context.Context, params exchanger.ConvertCurrencyParams) ([]exchanger.CurrencyRate, error) {
	if params.ConvertFrom == params.ConvertTo {
		var res []exchanger.CurrencyRate
		date := params.Start
		for date.Before(params.End.Add(time.Hour)) {
			res = append(res, exchanger.CurrencyRate{
				Base:     params.ConvertFrom,
				Rated:    params.ConvertTo,
				Sale:     1,
				Purchase: 1,
				Source:   "NBU",
				Date:     date,
			})

			date = date.Add(time.Hour * 24)
		}
	}

	if params.ConvertFrom == exchanger.UAH {
		return r.rater.FindRates(ctx, params.ConvertTo, params.Start, params.End)
	}

	rates, err := r.rater.FindRates(ctx, params.ConvertFrom, params.Start, params.End)
	if err != nil {
		return nil, fmt.Errorf("find rates: %w", err)
	}

	res := make([]exchanger.CurrencyRate, 0, len(rates))
	for _, rate := range rates {
		res = append(res, exchanger.CurrencyRate{
			Base:     params.ConvertFrom,
			Rated:    params.ConvertTo,
			Sale:     1 / rate.Sale,
			Purchase: 1 / rate.Purchase,
			Date:     rate.Date,
			Source:   rate.Source,
		})
	}

	return res, nil
}
