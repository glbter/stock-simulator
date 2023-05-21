package postgres

import (
	"context"
	"github.com/glbter/currency-ex/stocks"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTickerRepository_SaveSplits(t *testing.T) {
	type record struct {
		Date     time.Time `db:"date"`
		TickerID string    `db:"ticker_id"`
		Before   float64   `db:"before"`
		After    float64   `db:"after"`
	}

	now := time.Date(2013, time.April, 1, 0, 0, 0, 0, time.UTC)

	for tName, tCase := range map[string]struct {
		splits []stocks.SaveSplitParams
	}{
		"ok": {
			splits: []stocks.SaveSplitParams{
				{
					Date:     now,
					TickerID: "aad17418-6764-4ecd-90ed-bb1d7091edcc",
					Before:   1,
					After:    2,
				},
				{
					Date:     now.Add(time.Hour * 24 * 3),
					TickerID: "3439d561-b4db-4455-aff9-da2119573574",
					Before:   1,
					After:    2,
				},
			},
		},
		"empty": {},
	} {
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)
			fx.mustExecSQLFixture(t, testDB, "save_tickers.sql")

			require.NoError(t, TickerRepository{}.
				SaveSplits(context.Background(), testDB, tCase.splits))

			var res []record
			require.NoError(t,
				testDB.Select(
					context.Background(),
					&res,
					"select date, ticker_id, before, after from split",
				),
			)

			require.Len(t, res, len(tCase.splits))
			for i := range res {
				require.Equal(t, tCase.splits[i].Date, res[i].Date)
				require.Equal(t, tCase.splits[i].TickerID, res[i].TickerID)
				require.Equal(t, tCase.splits[i].Before, res[i].Before)
				require.Equal(t, tCase.splits[i].After, res[i].After)
			}
		})
	}
}

func TestTickerRepository_SaveDaily(t *testing.T) {
	type record struct {
		ID       string    `db:"id"`
		Date     time.Time `db:"date"`
		TickerID string    `db:"ticker_id"`
		High     *float64  `db:"high"`
		Low      *float64  `db:"low"`
		Open     *float64  `db:"open"`
		Close    *float64  `db:"close"`
		Volume   *float64  `db:"volume"`
	}

	now := time.Date(2013, time.April, 1, 0, 0, 0, 0, time.UTC)

	for tName, tCase := range map[string]struct {
		daily []stocks.SaveDailyParams
	}{
		"ok": {
			daily: []stocks.SaveDailyParams{
				{
					Date:     now,
					TickerID: "aad17418-6764-4ecd-90ed-bb1d7091edcc",
					High:     pointerFloat64(1),
					Low:      pointerFloat64(2),
					Open:     pointerFloat64(3),
					Close:    pointerFloat64(4),
					Volume:   pointerFloat64(5),
				},
				{
					Date:     now.Add(time.Hour * 24 * 3),
					TickerID: "3439d561-b4db-4455-aff9-da2119573574",
				},
			},
		},
		"empty": {},
	} {
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)
			fx.mustExecSQLFixture(t, testDB, "save_tickers.sql")

			require.NoError(t, TickerRepository{}.
				SaveDaily(context.Background(), testDB, tCase.daily))

			var res []record
			require.NoError(t,
				testDB.Select(
					context.Background(),
					&res,
					"select id, ticker_id, date, high, low, open, close, volume from stock_daily",
				),
			)

			require.Len(t, res, len(tCase.daily))
			for i := range res {
				require.NotEqual(t, "", res[i].ID)
				require.Equal(t, tCase.daily[i].Date, res[i].Date)
				require.Equal(t, tCase.daily[i].TickerID, res[i].TickerID)
				require.Equal(t, tCase.daily[i].High, res[i].High)
				require.Equal(t, tCase.daily[i].Low, res[i].Low)
				require.Equal(t, tCase.daily[i].Open, res[i].Open)
				require.Equal(t, tCase.daily[i].Close, res[i].Close)
				require.Equal(t, tCase.daily[i].Volume, res[i].Volume)
			}
		})
	}
}

func TestTickerRepository_SaveTicker(t *testing.T) {
	type record struct {
		ID   string  `db:"id"`
		Name string  `db:"name"`
		Desc *string `db:"description"`
	}

	for tName, tCase := range map[string]struct {
		tickers []stocks.SaveTickerParams
	}{
		"ok": {
			tickers: []stocks.SaveTickerParams{
				{
					Name: "ticker1",
				},
				{
					Name:        "ticker2",
					Description: pointerString("ticker"),
				},
			},
		},
		"empty": {},
	} {
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)

			require.NoError(t, TickerRepository{}.
				SaveTicker(context.Background(), testDB, tCase.tickers))

			var res []record
			require.NoError(t,
				testDB.Select(
					context.Background(),
					&res,
					"select id, description, name from ticker",
				),
			)

			require.Len(t, res, len(tCase.tickers))
			for i := range res {
				require.NotEqual(t, "", res[i].ID)
				require.Equal(t, tCase.tickers[i].Description, res[i].Desc)
				require.Equal(t, tCase.tickers[i].Name, res[i].Name)
			}
		})
	}
}
