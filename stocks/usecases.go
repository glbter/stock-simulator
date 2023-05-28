package stocks

import (
	"context"
	"github.com/glbter/currency-ex/currency/exchanger"
)

type ExchangeParams struct {
	ConverFrom exchanger.Currency
	ConvertTo  exchanger.Currency
}

type TickerUsecases interface {
	QueryLatestDaily(ctx context.Context, f QueryDailyFilter, ep ExchangeParams) ([]TickerWithData, error)
	QueryTickerDailyGraph(ctx context.Context, f QueryDailyGraphParams, ep ExchangeParams) ([]StockDailyData, error)
}

type PortfolioUsecases interface {
	CountPortfolio(ctx context.Context, userID string, ep ExchangeParams) (PortfolioState, error)
	TradeTickers(ctx context.Context, p TradeTickerParams) error
}
