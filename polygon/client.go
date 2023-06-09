package polygon

import (
	"context"
	"errors"
	"fmt"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	"github.com/glbter/currency-ex/stocks"
	"github.com/glbter/currency-ex/stocks/repository/postgres"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"time"
)

func NewDailyDataProcessor(
	pc *polygon.Client,
	s sqlc.SelectExecutor,
	repo postgres.TickerRepository,
	config Config,
) DailyDataProcessor {
	return DailyDataProcessor{
		pc:                 pc,
		s:                  s,
		repo:               repo,
		interval:           config.Timespan,
		intervalMultiplier: config.Multiplier,
		from:               config.From,
		to:                 config.To,
	}
}

type Config struct {
	Multiplier int
	Timespan   string
	From       time.Time
	To         time.Time
}

type DailyDataProcessor struct {
	pc   *polygon.Client
	s    sqlc.SelectExecutor
	repo postgres.TickerRepository

	intervalMultiplier int
	interval           string
	from               time.Time
	to                 time.Time
}

func (p DailyDataProcessor) Process(ctx context.Context, tickers []string) error {
	for _, ticker := range tickers {
		if err := p.SaveTickerData(ctx, ticker); err != nil {
			return fmt.Errorf("ticker %v: %w", ticker, err)
		}
	}

	return nil
}

func (p DailyDataProcessor) WithNewConfig(c Config) DailyDataProcessor {
	m := p.intervalMultiplier
	if c.Multiplier != 0 {
		m = c.Multiplier
	}

	i := p.interval
	if c.Timespan != "" {
		i = c.Timespan
	}

	f := p.from
	if !c.From.IsZero() {
		f = c.From
	}

	t := p.to
	if !c.To.IsZero() {
		t = c.To
	}

	return DailyDataProcessor{
		pc:                 p.pc,
		s:                  p.s,
		repo:               p.repo,
		interval:           i,
		intervalMultiplier: m,
		from:               f,
		to:                 t,
	}
}

func (p DailyDataProcessor) SaveTickerData(ctx context.Context, ticker string) error {
	tickers, err := p.repo.QueryTickers(ctx, p.s, stocks.QueryTickersFilters{Tickers: []string{ticker}})
	if err != nil {
		return err
	}

	id := tickers[0].ID

	dailyData, err := p.FetchDailyData(ctx, ticker, id)
	if err != nil {
		return err
	}

	return p.repo.SaveDaily(ctx, p.s, dailyData)
}

func (p DailyDataProcessor) FetchDailyData(ctx context.Context, ticker, tickerID string) ([]stocks.SaveDailyParams, error) {
	adj := true
	iter := p.pc.ListAggs(ctx, &models.ListAggsParams{
		Ticker:     ticker,
		Multiplier: p.intervalMultiplier,
		Timespan:   models.Timespan(p.interval),
		Adjusted:   &adj,
		From:       models.Millis(p.from),
		To:         models.Millis(p.to),
	})
	if iter == nil {
		return nil, errors.New("iterator is nil")
	}

	var dailyData []stocks.SaveDailyParams

	for iter.Next() {
		item := iter.Item()

		h := item.High
		l := item.Low
		o := item.Open
		c := item.Close
		v := item.Volume

		dailyData = append(
			dailyData,
			stocks.SaveDailyParams{
				TickerID: tickerID,
				Date:     time.Time(item.Timestamp),
				High:     &h,
				Low:      &l,
				Open:     &o,
				Close:    &c,
				Volume:   &v,
			},
		)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return dailyData, nil
}
