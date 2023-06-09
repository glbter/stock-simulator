package usecases

import (
	"context"
	"github.com/glbter/currency-ex/currency/exchanger"
	exchmock "github.com/glbter/currency-ex/currency/exchanger/mock"
	dbmock "github.com/glbter/currency-ex/pkg/sql/mock"
	"github.com/glbter/currency-ex/stocks"
	"github.com/glbter/currency-ex/stocks/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTickerInteractor_QueryLatestDaily(t *testing.T) {
	for tName, tCase := range map[string]struct {
		inputFilter   stocks.QueryDailyFilter
		inputExchange stocks.ExchangeParams

		expRes       []stocks.TickerWithData
		repoRes      []stocks.TickerWithData
		exchangeRate []exchanger.CurrencyRate
	}{
		"ok": {
			inputExchange: stocks.ExchangeParams{
				ConvertFrom: exchanger.USD,
				ConvertTo:   exchanger.UAH,
			},
			repoRes: []stocks.TickerWithData{
				{
					Ticker: stocks.Ticker{
						ID:   "1",
						Name: "aapl",
					},
					Volume: 3,
					High:   2,
					Low:    2,
					Open:   2,
					Close:  2,
				},
			},
			exchangeRate: []exchanger.CurrencyRate{
				{
					Base:     exchanger.USD,
					Rated:    exchanger.UAH,
					Sale:     2,
					Purchase: 3,
				},
			},
			expRes: []stocks.TickerWithData{
				{
					Ticker: stocks.Ticker{
						ID:   "1",
						Name: "aapl",
					},
					Volume: 3,
					High:   6,
					Low:    6,
					Open:   6,
					Close:  6,
				},
			},
		},
	} {
		tCase := tCase
		t.Run(tName, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := dbmock.NewMockDB(ctrl)
			exch := exchmock.NewMockAllCurrencyRater(ctrl)
			tr := mock.NewMockTickerRepository(ctrl)

			tr.EXPECT().
				QueryLatestDaily(gomock.Any(), gomock.Any(), tCase.inputFilter).
				Return(tCase.repoRes, nil)

			if tCase.exchangeRate != nil {
				exch.EXPECT().
					FindRates(gomock.Any(), gomock.Any()).
					Do(func(_ context.Context, params exchanger.ConvertCurrencyParams) {
						require.Equal(t, params.ConvertTo, tCase.inputExchange.ConvertTo)
						require.Equal(t, params.ConvertFrom, tCase.inputExchange.ConvertFrom)
					}).
					Return(tCase.exchangeRate, nil)
			}

			inter := NewTickerInteractor(db, tr, exch)
			res, err := inter.QueryLatestDaily(context.Background(), tCase.inputFilter, tCase.inputExchange)

			require.NoError(t, err)
			require.Equal(t, tCase.expRes, res)
		})
	}
}
