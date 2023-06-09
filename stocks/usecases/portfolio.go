package usecases

import (
	"context"
	"fmt"
	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/glbter/currency-ex/pkg/serrors"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	"github.com/glbter/currency-ex/stocks"
	"math/rand"
	"time"
)

type PortfolioInteractor struct {
	db            sqlc.DB
	repo          stocks.PortfolioRepository
	exchangeRater exchanger.AllCurrencyRater
	tickerRepo    stocks.TickerRepository
}

func NewPortfolioInteractor(
	db sqlc.DB,
	repo stocks.PortfolioRepository,
	exchangeRater exchanger.AllCurrencyRater,
	tickerRepo stocks.TickerRepository,
) PortfolioInteractor {
	return PortfolioInteractor{
		db:            db,
		repo:          repo,
		exchangeRater: exchangeRater,
		tickerRepo:    tickerRepo,
	}
}

func (i PortfolioInteractor) TradeTickers(ctx context.Context, p stocks.TradeTickerParams) error {
	// change to ticker name maybe?
	amount, err := i.repo.CountTickerAmount(ctx, i.db, stocks.CountTickerAmountParams{
		UserID:    p.UserID,
		TickerIDs: []string{p.TickerID},
	})
	if err != nil {
		return fmt.Errorf("count amount: %w", err)
	}

	if p.Action == stocks.ACTION_SELL {
		if len(amount) != 1 || amount[0].Amount-p.Amount < 0 {
			return fmt.Errorf("%w: %v", serrors.ErrBadInput, "not enough tickers")
		}
	}

	dailies, err := i.tickerRepo.QueryLatestDaily(ctx, i.db, stocks.QueryDailyFilter{
		TickerIDs: []string{p.TickerID},
	})
	if err != nil {
		return fmt.Errorf("query latest daily: %w", err)
	}

	if len(dailies) != 1 {
		return fmt.Errorf("more than one daily: %w", err)
	}

	daily := dailies[0]
	price := simulatePrice(daily.Low, daily.High)

	// TODO: add trade by price limit
	if err := i.repo.TradeTickers(ctx, i.db, stocks.TradeTickerParams{
		TickerID: p.TickerID,
		UserID:   p.UserID,
		Amount:   p.Amount,
		PriceUSD: price,
		Action:   p.Action,
	}); err != nil {
		return fmt.Errorf("trade ticker: %w", err)
	}

	return nil
}

func simulatePrice(low, high float64) float64 {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	mean := (high + low) / 2
	stdDev := (high - low) / 6
	return r.NormFloat64()*stdDev + mean
}

func (i PortfolioInteractor) CountPortfolio(ctx context.Context, userID string, ep stocks.ExchangeParams) (stocks.PortfolioState, error) {
	state, err := i.repo.CountPortfolio(ctx, i.db, userID)
	if err != nil {
		return stocks.PortfolioState{}, fmt.Errorf("count portfolio: %w", err)
	}

	if len(state.All) == 0 {
		return state, nil
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	rates, err := i.exchangeRater.FindRates(ctx, exchanger.ConvertCurrencyParams{
		ConvertFrom: ep.ConvertFrom,
		ConvertTo:   ep.ConvertTo,
		Start:       today,
		End:         today,
	})
	if err != nil {
		return stocks.PortfolioState{}, fmt.Errorf("find rate: %w", err)
	}

	if len(rates) != 1 {
		return stocks.PortfolioState{}, fmt.Errorf("find rate: %w", err)
	}

	rate := rates[0]
	res := stocks.PortfolioState{
		Total: convertRateOfPortfolioState(state.Total, rate),
		All:   make([]stocks.PortfolioTickerState, 0, len(state.All)),
	}

	for _, ps := range state.All {
		res.All = append(res.All, convertRateOfPortfolioState(ps, rate))
	}

	return res, nil
}

func convertRateOfPortfolioState(ps stocks.PortfolioTickerState, rate exchanger.CurrencyRate) stocks.PortfolioTickerState {
	return stocks.PortfolioTickerState{
		TickerID:    ps.TickerID,
		Name:        ps.Name,
		Description: ps.Description,
		Amount:      ps.Amount,
		High:        ps.High * rate.Sale,
		Low:         ps.Low * rate.Sale,
		Open:        ps.Open * rate.Sale,
		Close:       ps.Close * rate.Sale,
	}
}
