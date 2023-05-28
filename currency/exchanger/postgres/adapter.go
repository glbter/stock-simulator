package postgres

import (
	"context"
	"github.com/glbter/currency-ex/currency/exchanger"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	"time"
)

func NewCurryingExchangeAdapter(s sqlc.Selector, r Rater, source exchanger.ExchangeSource) Adapter {
	return Adapter{r: r, s: s, source: source}
}

type Adapter struct {
	r      Rater
	s      sqlc.Selector
	source exchanger.ExchangeSource
}

func (a Adapter) FindRates(
	ctx context.Context,
	c exchanger.Currency,
	start time.Time,
	end time.Time,
) ([]exchanger.CurrencyRate, error) {
	return a.r.FindRates(ctx, a.s,
		FindExchangeRateParams{
			Currency: c,
			Start:    &start,
			End:      &end,
			Source:   a.source,
		},
	)
}
