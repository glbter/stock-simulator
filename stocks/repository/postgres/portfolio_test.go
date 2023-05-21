package postgres

import (
	"context"
	"github.com/glbter/currency-ex/stocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPortfolioRepository_CountPortfolio(t *testing.T) {
	defer mustTruncateTables(t, testDB)
	fx.mustExecSQLFixture(t, testDB, "data.sql")

	for tName, tCase := range map[string]struct {
		userID string

		expected stocks.PortfolioState
	}{
		"user_has_tickers": {
			userID: "4ffdaa1c-9c25-4a79-a3a6-cf47ba361728",
			expected: stocks.PortfolioState{
				Total: stocks.PortfolioTickerState{
					TickerID:    "00000000-0000-0000-0000-000000000000",
					Amount:      7,
					Name:        "TOTAL",
					Description: "TOTAL",
					High:        21,
					Low:         14,
					Open:        28,
					Close:       35,
				},
				All: []stocks.PortfolioTickerState{
					{
						TickerID: "aad17418-6764-4ecd-90ed-bb1d7091edcc",
						Amount:   6,
						Name:     "AAPL",
						High:     18,
						Low:      12,
						Open:     24,
						Close:    30,
					},
					{
						TickerID: "3439d561-b4db-4455-aff9-da2119573574",
						Amount:   1,
						Name:     "AAPL2",
						High:     3,
						Low:      2,
						Open:     4,
						Close:    5,
					},
				},
			},
		},
		"user_has_no_tickers": {
			userID: "00000000-0000-0000-0000-000000000000",
			expected: stocks.PortfolioState{
				Total: stocks.PortfolioTickerState{
					TickerID:    "00000000-0000-0000-0000-000000000000",
					Name:        "TOTAL",
					Description: "TOTAL",
				},
				All: []stocks.PortfolioTickerState{},
			},
		},
	} {
		t.Run(tName, func(t *testing.T) {
			res, err := PortfolioRepository{}.
				CountPortfolio(context.Background(), testDB, tCase.userID)

			require.NoError(t, err)
			require.Equal(t, tCase.expected, res)
		})
	}
}

func TestPortfolioRepository_TradeTickers(t *testing.T) {
	type record struct {
		ID         string  `db:"id"`
		InvestorID string  `db:"investor_id"`
		TickerID   string  `db:"ticker_id"`
		Amount     float64 `db:"amount"`
		Price      float64 `db:"price_usd"`
		Action     string  `db:"action"`
	}

	for tName, tCase := range map[string]struct {
		params stocks.TradeTickerParams

		expected record
	}{
		"buy": {
			params: stocks.TradeTickerParams{
				UserID:   "9b8b404d-9d33-4c23-9ca6-2d98ce399a00",
				TickerID: "fdd2c4f2-cb0a-4139-9ae4-f0a8e281aadb",
				Amount:   65,
				Action:   stocks.ACTION_BUY,
				PriceUSD: 40,
			},
			expected: record{
				InvestorID: "9b8b404d-9d33-4c23-9ca6-2d98ce399a00",
				TickerID:   "fdd2c4f2-cb0a-4139-9ae4-f0a8e281aadb",
				Amount:     65,
				Action:     string(stocks.ACTION_BUY),
				Price:      40,
			},
		},
		"sell": {
			params: stocks.TradeTickerParams{
				UserID:   "9b8b404d-9d33-4c23-9ca6-2d98ce399a00",
				TickerID: "fdd2c4f2-cb0a-4139-9ae4-f0a8e281aadb",
				Amount:   65,
				Action:   stocks.ACTION_SELL,
				PriceUSD: 40,
			},
			expected: record{
				InvestorID: "9b8b404d-9d33-4c23-9ca6-2d98ce399a00",
				TickerID:   "fdd2c4f2-cb0a-4139-9ae4-f0a8e281aadb",
				Amount:     65,
				Action:     string(stocks.ACTION_SELL),
				Price:      40,
			},
		},
	} {
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)

			require.NoError(t,
				PortfolioRepository{}.
					TradeTickers(context.Background(), testDB, tCase.params),
			)

			var res []record
			require.NoError(t,
				testDB.Select(
					context.Background(),
					&res,
					"select id, investor_id, ticker_id, amount, price_usd, action from portfolio_record",
				),
			)

			require.Len(t, res, 1)
			require.NotEqual(t, "", res[0].ID)
			require.Equal(t, tCase.expected.InvestorID, res[0].InvestorID)
			require.Equal(t, tCase.expected.TickerID, res[0].TickerID)
			require.Equal(t, tCase.expected.Amount, res[0].Amount)
			require.Equal(t, tCase.expected.Price, res[0].Price)
			require.Equal(t, tCase.expected.Action, res[0].Action)
		})
	}
}
