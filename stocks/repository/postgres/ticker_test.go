package postgres

import (
	"context"
	"encoding/json"
	"fmt"
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
		tCase := tCase
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)
			fx.MustExecSQLFixture(t, testDB, "save_tickers.sql")

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
		tCase := tCase
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)
			fx.MustExecSQLFixture(t, testDB, "save_tickers.sql")

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
		tCase := tCase
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

func TestTickerRepository_QueryLatestDaily(t *testing.T) {
	for tName, tCase := range map[string]struct {
		filter stocks.QueryDailyFilter
		exp    []stocks.TickerWithData
	}{
		"one_ticker": {
			filter: stocks.QueryDailyFilter{
				TickerIDs: []string{"aad17418-6764-4ecd-90ed-bb1d7091edcc"},
			},
			exp: []stocks.TickerWithData{
				{
					Ticker: stocks.Ticker{
						ID:   "aad17418-6764-4ecd-90ed-bb1d7091edcc",
						Name: "AAPL",
					},
					High:     3,
					Low:      2,
					Open:     4,
					Close:    5,
					DataDate: time.Date(2010, 01, 20, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		"two_tickers": {
			filter: stocks.QueryDailyFilter{
				TickerIDs: []string{"aad17418-6764-4ecd-90ed-bb1d7091edcc", "3439d561-b4db-4455-aff9-da2119573574"},
			},
			exp: []stocks.TickerWithData{
				{
					Ticker: stocks.Ticker{
						ID:   "aad17418-6764-4ecd-90ed-bb1d7091edcc",
						Name: "AAPL",
					},
					High:     3,
					Low:      2,
					Open:     4,
					Close:    5,
					DataDate: time.Date(2010, 01, 20, 0, 0, 0, 0, time.UTC),
				},
				{
					Ticker: stocks.Ticker{
						ID:   "3439d561-b4db-4455-aff9-da2119573574",
						Name: "AAPL2",
					},
					High:     3,
					Low:      2,
					Open:     4,
					Close:    5,
					DataDate: time.Date(2010, 01, 20, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		"by_name": {
			filter: stocks.QueryDailyFilter{
				Tickers: []string{"AAPL"},
			},
			exp: []stocks.TickerWithData{
				{
					Ticker: stocks.Ticker{
						ID:   "aad17418-6764-4ecd-90ed-bb1d7091edcc",
						Name: "AAPL",
					},
					High:     3,
					Low:      2,
					Open:     4,
					Close:    5,
					DataDate: time.Date(2010, 01, 20, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		"not_existing_ticker": {
			filter: stocks.QueryDailyFilter{
				TickerIDs: []string{"38dcdca9-a2fe-4b46-8f10-aa1bd0fd26e7"},
			},
			exp: []stocks.TickerWithData{},
		},
		"all": {
			filter: stocks.QueryDailyFilter{
				TickerIDs: []string{},
			},
			exp: []stocks.TickerWithData{
				{
					Ticker: stocks.Ticker{
						ID:   "aad17418-6764-4ecd-90ed-bb1d7091edcc",
						Name: "AAPL",
					},
					High:     3,
					Low:      2,
					Open:     4,
					Close:    5,
					DataDate: time.Date(2010, 01, 20, 0, 0, 0, 0, time.UTC),
				},
				{
					Ticker: stocks.Ticker{
						ID:   "3439d561-b4db-4455-aff9-da2119573574",
						Name: "AAPL2",
					},
					High:     3,
					Low:      2,
					Open:     4,
					Close:    5,
					DataDate: time.Date(2010, 01, 20, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	} {
		tCase := tCase
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)

			fx.MustExecSQLFixture(t, testDB, "data2.sql")
			d, err := TickerRepository{}.
				QueryLatestDaily(context.Background(), testDB, tCase.filter)

			require.NoError(t, err)
			require.ElementsMatch(t, tCase.exp, d)
		})
	}
}

func TestTickerRepository_QueryTickers(t *testing.T) {
	for tName, tCase := range map[string]struct {
		filter   stocks.QueryTickersFilters
		expected []stocks.Ticker
	}{
		"all": {
			filter: stocks.QueryTickersFilters{},
			expected: []stocks.Ticker{
				{
					ID:   "aad17418-6764-4ecd-90ed-bb1d7091edcc",
					Name: "AAPL",
				},
				{
					ID:   "3439d561-b4db-4455-aff9-da2119573574",
					Name: "AAPL2",
				},
			},
		},
		"ticker_id": {
			filter: stocks.QueryTickersFilters{
				IDs: []string{"aad17418-6764-4ecd-90ed-bb1d7091edcc"},
			},
			expected: []stocks.Ticker{
				{
					ID:   "aad17418-6764-4ecd-90ed-bb1d7091edcc",
					Name: "AAPL",
				},
			},
		},
		"ticker_name": {
			filter: stocks.QueryTickersFilters{
				Tickers: []string{"AAPL2"},
			},
			expected: []stocks.Ticker{
				{
					ID:   "3439d561-b4db-4455-aff9-da2119573574",
					Name: "AAPL2",
				},
			},
		},
		"ticker_name_id": {
			filter: stocks.QueryTickersFilters{
				IDs:     []string{"aad17418-6764-4ecd-90ed-bb1d7091edcc"},
				Tickers: []string{"AAPL2"},
			},
			expected: []stocks.Ticker{},
		},
	} {
		tCase := tCase
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)

			fx.MustExecSQLFixture(t, testDB, "data.sql")

			d, err := TickerRepository{}.
				QueryTickers(context.Background(), testDB, tCase.filter)

			require.NoError(t, err)
			require.ElementsMatch(t, tCase.expected, d)
		})
	}
}

func TestTickerRepository_QueryTickerDailyGraph(t *testing.T) {
	for tName, tCase := range map[string]struct {
		filter     stocks.QueryDailyGraphParams
		expFixture string
	}{
		"all_for_ticker": {
			filter: stocks.QueryDailyGraphParams{
				TickerIDs: []string{"3439d561-b4db-4455-aff9-da2119573574"},
			},
			expFixture: "all_for_ticker.json",
		},
		"all": {
			filter:     stocks.QueryDailyGraphParams{},
			expFixture: "all.json",
		},
		"before": {
			filter: stocks.QueryDailyGraphParams{
				TickerIDs:     []string{"3439d561-b4db-4455-aff9-da2119573574"},
				BeforeDateInc: pointerTime(time.Date(2010, 01, 24, 0, 0, 0, 0, time.UTC)), //2010-01-24T00:00:00Z
			},
			expFixture: "before.json",
		},
		"after": {
			filter: stocks.QueryDailyGraphParams{
				TickerIDs:    []string{"3439d561-b4db-4455-aff9-da2119573574"},
				AfterDateInc: pointerTime(time.Date(2010, 01, 21, 0, 0, 0, 0, time.UTC)), //2010-01-21T00:00:00Z
			},
			expFixture: "after.json",
		},
		"before_after": {
			filter: stocks.QueryDailyGraphParams{
				TickerIDs:     []string{"3439d561-b4db-4455-aff9-da2119573574"},
				BeforeDateInc: pointerTime(time.Date(2010, 01, 24, 0, 0, 0, 0, time.UTC)), //2010-01-24T00:00:00Z
				AfterDateInc:  pointerTime(time.Date(2010, 01, 21, 0, 0, 0, 0, time.UTC)), //2010-01-24T00:00:00Z
			},
			expFixture: "before_after.json",
		},
	} {
		tfName := t.Name()
		tCase := tCase
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)

			fx.MustExecSQLFixture(t, testDB, "query_daily_graph.sql")

			fixture := fx.MustLoadStringFixture(t, fmt.Sprintf("%v/%v", tfName, tCase.expFixture))

			d, err := TickerRepository{}.
				QueryTickerDailyGraph(context.Background(), testDB, tCase.filter)

			require.NoError(t, err)

			res, err := json.Marshal(d)
			require.NoError(t, err)

			require.JSONEq(t, fixture, string(res))
		})
	}
}
