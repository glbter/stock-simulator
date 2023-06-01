package usecases

import (
	"context"
	"fmt"
	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/glbter/currency-ex/currency/exchanger/factory"
	"github.com/glbter/currency-ex/pkg/serrors"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	"github.com/glbter/currency-ex/stocks"
)

type TickerInteractor struct {
	db            sqlc.DB
	repo          stocks.TickerRepository
	exchangeRater factory.AllCurrencyRater
}

func NewTickerInteractor(
	db sqlc.DB,
	repo stocks.TickerRepository,
	exchangeRater factory.AllCurrencyRater,
) TickerInteractor {
	return TickerInteractor{
		db:            db,
		exchangeRater: exchangeRater,
		repo:          repo,
	}
}

func (i TickerInteractor) QueryTickers(ctx context.Context, f stocks.QueryTickersFilters) ([]stocks.Ticker, error) {
	return i.repo.QueryTickers(ctx, i.db, f)
}

func (i TickerInteractor) QueryLatestDaily(
	ctx context.Context,
	f stocks.QueryDailyFilter,
	ep stocks.ExchangeParams,
) ([]stocks.TickerWithData, error) {
	tickerDaily, err := i.repo.QueryLatestDaily(ctx, i.db, f)
	if err != nil {
		return nil, fmt.Errorf("query daily: %w", err)
	}

	//return tickerDaily, nil

	res := make([]stocks.TickerWithData, 0, len(tickerDaily))
	for _, d := range tickerDaily {
		rates, err := i.exchangeRater.FindRates(ctx, exchanger.ConvertCurrencyParams{
			ConvertFrom: ep.ConverFrom,
			ConvertTo:   ep.ConvertTo,
			Start:       d.DataDate,
			End:         d.DataDate,
		})
		if err != nil {
			return nil, fmt.Errorf("find rate: %w", err)
		}

		// TODO: fix
		if len(rates) != 1 {
			continue
		}

		rate := rates[0]

		res = append(res, stocks.TickerWithData{
			Ticker:   d.Ticker,
			DataDate: d.DataDate,
			Volume:   d.Volume,
			High:     d.High * rate.Purchase,
			Low:      d.Low * rate.Purchase,
			Open:     d.Open * rate.Purchase,
			Close:    d.Close * rate.Purchase,
		})
	}

	return res, nil
}
func (i TickerInteractor) QueryTickerDailyGraph(ctx context.Context, f stocks.QueryDailyGraphParams, ep stocks.ExchangeParams) ([]stocks.StockDailyData, error) {
	data, err := i.repo.QueryTickerDailyGraph(ctx, i.db, f)
	if err != nil {
		return nil, fmt.Errorf("query daily graph: %w", err)
	}

	//return data, nil

	rates, err := i.exchangeRater.FindRates(ctx, exchanger.ConvertCurrencyParams{
		ConvertFrom: ep.ConverFrom,
		ConvertTo:   ep.ConvertTo,
		Start:       data[0].Date,
		End:         data[len(data)-1].Date,
	})
	if err != nil {
		return nil, fmt.Errorf("find rates: %w", err)
	}

	res := make([]stocks.StockDailyData, 0, len(rates))

	var g int
	var j int
	if len(data) == 0 {
		return nil, fmt.Errorf("no graph found: %w", serrors.ErrNotFound)
	}
	if len(rates) == 0 {
		return nil, fmt.Errorf("no rates found: %w", serrors.ErrNotFound)
	}

	for n := 0; n < len(data) || n < len(rates); n++ {
		if data[g].Date.Equal(rates[j].Date) {
			res = append(res, data[g].MultiplyPrice(rates[j].Purchase))
		}

		if data[g].Date.Before(rates[j].Date) {
			j++
			continue
		}

		if data[g].Date.After(rates[j].Date) {
			//TODO: fix
			res = append(res, data[g].MultiplyPrice(rates[j].Purchase))
		}
	}

	return res, nil
}
