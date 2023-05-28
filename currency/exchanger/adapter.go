package exchanger

import (
	"context"
	"fmt"
	"time"
)

type AllCurrencyRater struct {
	rater CurrencySeriesRater
}

func NewAllCurrencyRater(rater CurrencySeriesRater) AllCurrencyRater {
	return AllCurrencyRater{
		rater: rater,
	}
}

type ConvertCurrencyParams struct {
	ConvertFrom Currency
	ConvertTo   Currency
	Start       time.Time
	End         time.Time
}

func (r AllCurrencyRater) FindRates(ctx context.Context, params ConvertCurrencyParams) ([]CurrencyRate, error) {
	if params.ConvertFrom == params.ConvertTo {
		var res []CurrencyRate
		date := params.Start
		for date.Before(params.End.Add(time.Hour)) {
			res = append(res, CurrencyRate{
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

	if params.ConvertFrom == UAH {
		return r.rater.FindRates(ctx, params.ConvertTo, params.Start, params.End)
	}

	rates, err := r.rater.FindRates(ctx, params.ConvertFrom, params.Start, params.End)
	if err != nil {
		return nil, fmt.Errorf("find rates: %w", err)
	}

	res := make([]CurrencyRate, 0, len(rates))
	for _, rate := range rates {
		res = append(res, CurrencyRate{
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
