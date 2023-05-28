package stocks

import (
	"context"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	"time"
)

type PortfolioRepository interface {
	TradeTickers(ctx context.Context, e sqlc.Executor, p TradeTickerParams) error
	CountTickerAmount(ctx context.Context, s sqlc.Selector, p CountTickerAmountParams) ([]PortfolioTickerAmount, error)
	CountPortfolio(ctx context.Context, s sqlc.Selector, userID string) (PortfolioState, error)
}

type TradeTickerParams struct {
	UserID   string
	TickerID string
	Amount   float64
	Action   PortfolioAction
	PriceUSD float64
}

type CountTickerAmountParams struct {
	UserID    string
	TickerIDs []string
}

type TickerRepository interface {
	SaveSplits(ctx context.Context, e sqlc.Executor, ps []SaveSplitParams) error
	SaveDaily(ctx context.Context, e sqlc.Executor, ps []SaveDailyParams) error
	SaveTicker(ctx context.Context, e sqlc.Executor, ps []SaveTickerParams) error

	QueryTickers(ctx context.Context, s sqlc.Selector, f QueryTickersFilters) ([]Ticker, error)
	QueryLatestDaily(ctx context.Context, s sqlc.Selector, f QueryDailyFilter) ([]TickerWithData, error)
	QueryTickerDailyGraph(ctx context.Context, s sqlc.Selector, f QueryDailyGraphParams) ([]StockDailyData, error)
}

type SaveSplitParams struct {
	Date     time.Time
	TickerID string
	Before   float64
	After    float64
}

type SaveDailyParams struct {
	TickerID string
	Date     time.Time
	High     *float64
	Low      *float64
	Open     *float64
	Close    *float64
	Volume   *float64
}

type SaveTickerParams struct {
	Name        string
	Description *string
}

type QueryTickersFilters struct {
	IDs     []string
	Tickers []string
}

type QueryDailyFilter struct {
	TickerIDs []string
}

type QueryDailyGraphParams struct {
	TickerIDs     []string
	BeforeDateInc *time.Time
	AfterDateInc  *time.Time
}
