package postgres

import (
	"context"
	"fmt"
	"github.com/glbter/currency-ex/currency/exchanger"
	sqlc "github.com/glbter/currency-ex/sql"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

const (
	currencyRateTable = "currency_rate"
)

var (
	currencyRateTableColumns = []string{"rate_date", "base_currency", "target_currency", "sale", "purchase", "source"}
)

type Rater struct{}

func (r Rater) SaveRate(ctx context.Context, executor sqlc.Executor, crs []exchanger.CurrencyRate) error {
	if len(crs) == 0 {
		return nil
	}

	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto(currencyRateTable).
		Cols(currencyRateTableColumns...)

	for _, cr := range crs {
		ib.Values(cr.Date, cr.Base, cr.Rated, cr.Sale, cr.Purchase, cr.Source)
	}

	q, args := ib.BuildWithFlavor(sqlbuilder.PostgreSQL)

	res, err := executor.Exec(ctx, q, args...)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if int(rows) != len(crs) {
		return fmt.Errorf("rows affected %v mismatch %v", rows, len(crs))
	}

	return nil
}

type rate struct {
	Date     time.Time `db:"rate_date"`
	Base     string    `db:"base_currency"`
	Target   string    `db:"target_currency"`
	Sale     float64   `db:"sale"`
	Purchase float64   `db:"purchase"`
	Source   string    `db:"source"`
}

func (r rate) ToCurrencyRate() (exchanger.CurrencyRate, error) {
	b, err := exchanger.ToCurrency(r.Base)
	if err != nil {
		return exchanger.CurrencyRate{}, fmt.Errorf("convert base: %w", err)
	}

	t, err := exchanger.ToCurrency(r.Target)
	if err != nil {
		return exchanger.CurrencyRate{}, fmt.Errorf("convert target: %w", err)
	}

	return exchanger.CurrencyRate{
		Base:     b,
		Rated:    t,
		Sale:     r.Sale,
		Purchase: r.Purchase,
		Date:     r.Date,
		Source:   exchanger.ExchangeSource(r.Source),
	}, nil
}

func (r Rater) FindRates(
	ctx context.Context,
	s sqlc.Selector,
	c exchanger.Currency,
	start time.Time,
	end time.Time,
	source exchanger.ExchangeSource,
) ([]exchanger.CurrencyRate, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(currencyRateTableColumns...).
		From(currencyRateTable).
		Where(sb.And(
			sb.Equal("target_currency", c),
			sb.Between("rate_date", start, end),
			sb.Equal("source", source),
		))

	q, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var rates []rate
	if err := s.Select(ctx, &rates, q, args...); err != nil {
		return nil, err
	}

	res := make([]exchanger.CurrencyRate, 0, len(rates))
	for _, rate := range rates {
		rr, err := rate.ToCurrencyRate()
		if err != nil {
			return nil, fmt.Errorf("convert rate for date %v: %w", rate.Date, err)
		}

		res = append(res, rr)
	}

	return res, nil
}
