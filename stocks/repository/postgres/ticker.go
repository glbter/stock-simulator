package postgres

import (
	"context"
	sqlc "github.com/glbter/currency-ex/sql"
	"github.com/glbter/currency-ex/stocks"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

func NewTickerRepository() TickerRepository {
	return TickerRepository{}
}

type TickerRepository struct {
}

var _ stocks.TickerRepository = &TickerRepository{}

func (TickerRepository) SaveSplits(ctx context.Context, e sqlc.Executor, ps []stocks.SaveSplitParams) error {
	if len(ps) == 0 {
		return nil
	}

	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("split").
		Cols("date", "ticker_id", "before", "after")

	for _, p := range ps {
		ib.Values(p.Date, p.TickerID, p.Before, p.After)
	}

	q, args := ib.BuildWithFlavor(sqlbuilder.PostgreSQL)

	if _, err := e.Exec(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

func (TickerRepository) SaveDaily(ctx context.Context, e sqlc.Executor, ps []stocks.SaveDailyParams) error {
	if len(ps) == 0 {
		return nil
	}

	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("stock_daily")
	ib.Cols("ticker_id", "date", "high", "low", "open", "close", "volume")

	for _, p := range ps {
		ib.Values(p.TickerID, p.Date, p.High, p.Low, p.Open, p.Close, p.Volume)
	}

	q, args := ib.BuildWithFlavor(sqlbuilder.PostgreSQL)
	if _, err := e.Exec(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

func (TickerRepository) SaveTicker(ctx context.Context, e sqlc.Executor, ps []stocks.SaveTickerParams) error {
	if len(ps) == 0 {
		return nil
	}

	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("ticker")
	ib.Cols("name", "description")

	for _, p := range ps {
		ib.Values(p.Name, p.Description)
	}

	q, args := ib.BuildWithFlavor(sqlbuilder.PostgreSQL)
	if _, err := e.Exec(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

type ticker struct {
	ID          string
	Name        string
	Destination *string
}

func (t ticker) toEntity() stocks.Ticker {
	var d string
	if t.Destination != nil {
		d = *t.Destination
	}

	return stocks.Ticker{
		ID:          t.ID,
		Name:        t.Name,
		Description: d,
	}
}

func (TickerRepository) QueryTickers(ctx context.Context, s sqlc.Selector, f stocks.QueryTickersFilters) ([]stocks.Ticker, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("id", "name", "description").
		From("ticker")

	if len(f.Tickers) != 0 {
		sb.Where(sb.In("name", sqlbuilder.List(f.Tickers)))
	}

	if len(f.IDs) != 0 {
		sb.Where(sb.In("id", sqlbuilder.List(f.IDs)))
	}

	q, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var tickers []ticker
	if err := s.Select(ctx, &tickers, q, args...); err != nil {
		return nil, err
	}

	res := make([]stocks.Ticker, 0, len(tickers))
	for _, ticker := range tickers {
		res = append(res, ticker.toEntity())
	}

	return res, nil
}

type tickerWithDaily struct {
	TickerID    string    `db:"ticker_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Date        time.Time `db:"date"`
	High        float64   `db:"high"`
	Low         float64   `db:"low"`
	Open        float64   `db:"open"`
	Close       float64   `db:"close"`
	Volume      float64   `db:"volume"`
}

func (d tickerWithDaily) toEntity() stocks.TickerWithData {
	return stocks.TickerWithData{
		Ticker: stocks.Ticker{
			ID:          d.TickerID,
			Name:        d.Name,
			Description: d.Description,
		},
		High:   d.High,
		Low:    d.Low,
		Open:   d.Open,
		Close:  d.Close,
		Volume: d.Volume,
	}
}

func (TickerRepository) QueryLatestDaily(ctx context.Context, s sqlc.Selector, f stocks.QueryDailyFilter) ([]stocks.TickerWithData, error) {
	sbLatestDaily := sqlbuilder.NewSelectBuilder()
	sbLatestDaily.
		Select(
			"ticker_id",
			"max(date) as date",
		).
		From("stock_daily").GroupBy("ticker_id")

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(
		"ticker.id",
		"ticker.name",
		"coalesce(ticker.description, '')",
		"sd.date",
		"coalesce(sd.high, 0)",
		"coalesce(sd.low, 0)",
		"coalesce(sd.open, 0)",
		"coalesce(sd.close, 0)",
		"coalesce(sd.volume, 0)",
	).
		From("ticker").
		Join(
			sb.BuilderAs(sbLatestDaily, "latest_daily"),
			"ticker.id = latest_daily.ticker_id",
		).
		Join("stock_daily",
			sb.And(
				"ticker.id = sd.ticker_id",
				"latest_daily.date = sd.date",
			),
		)

	if len(f.TickerIDs) > 0 {
		sb.Where(sb.In("ticker_id", sqlbuilder.List(f.TickerIDs)))
	}

	q, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var daily []tickerWithDaily
	if err := s.Select(ctx, &daily, q, args...); err != nil {
		return nil, err
	}

	res := make([]stocks.TickerWithData, 0, len(daily))
	for _, d := range daily {
		res = append(res, d.toEntity())
	}

	return res, nil

}

type daily struct {
	TickerID string    `db:"ticker_id"`
	Date     time.Time `db:"date"`
	High     float64   `db:"high"`
	Low      float64   `db:"low"`
	Open     float64   `db:"open"`
	Close    float64   `db:"close"`
	Volume   float64   `db:"volume"`
}

func (d daily) toEntity() stocks.StockDailyData {
	return stocks.StockDailyData{
		TickerID: d.TickerID,
		Date:     d.Date,
		High:     d.High,
		Low:      d.Low,
		Open:     d.Open,
		Close:    d.Close,
		Volume:   d.Volume,
	}
}

func (TickerRepository) QueryTickerDailyGraph(ctx context.Context, s sqlc.Selector, f stocks.QueryDailyGraphParams) ([]stocks.StockDailyData, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(
		"sd.ticker_id",
		"sd.date",
		"coalesce(sd.high, 0)",
		"coalesce(sd.low, 0)",
		"coalesce(sd.open, 0)",
		"coalesce(sd.close, 0)",
		"coalesce(sd.volume, 0)",
	).
		From(sb.As("stock_daily", "sd"))

	if len(f.TickerIDs) != 0 {
		sb.Where(sb.In("ticker_id", sqlbuilder.List(f.TickerIDs)))
	}

	if f.BeforeDateInc != nil {
		sb.Where(sb.LessEqualThan("sb.date", *f.BeforeDateInc))
	}

	if f.AfterDateInc != nil {
		sb.Where(sb.GreaterEqualThan("sb.date", *f.AfterDateInc))
	}

	q, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var daily []daily
	if err := s.Select(ctx, &daily, q, args...); err != nil {
		return nil, err
	}

	res := make([]stocks.StockDailyData, 0, len(daily))
	for _, d := range daily {
		res = append(res, d.toEntity())
	}

	return res, nil
}
