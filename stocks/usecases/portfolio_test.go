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

func TestPortfolioInteractor_TradeTickers(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//
	//r := mock.NewMockPortfolioRepository(ctrl)
	//db := dbmock.NewMockDB(ctrl)
	//exch := exchmock.NewMockAllCurrencyRater(ctrl)
	//tr := mock.NewMockTickerRepository(ctrl)
	//
	//r.EXPECT().
	//	CountTickerAmount(gomock.Any(), gomock.Any(), stocks.CountTickerAmountParams{UserID: "1", TickerIDs: []string{"1"}})
	//
	//r.EXPECT().
	//	TradeTickers(gomock.Any(), gomock.Any(),  stocks.TradeTickerParams{UserID: "1", TickerID: "1", Amount: 1, Action: stocks.ACTION_BUY}).
	//	Return(nil)
	//
	//inter := NewPortfolioInteractor(db, r, exch, tr)
	////inter := NewTickerInteractor(db, r)
	//inter.TradeTickers(context.Background(), stocks.TradeTickerParams{UserID: "1", TickerID: "1", Amount: 1, Action: stocks.ACTION_BUY})
	////inter.CountPortfolio(context.Background(), "1", stocks.ExchangeParams{
	////	ConvertFrom: exchanger.USD,
	////	ConvertTo:  exchanger.UAH,
	////})

}

func TestPortfolioInteractor_CountPortfolio(t *testing.T) {
	for tName, tCase := range map[string]struct {
		input        stocks.ExchangeParams
		expRes       stocks.PortfolioState
		repoRes      stocks.PortfolioState
		exchangeRate []exchanger.CurrencyRate
	}{
		"ok": {
			input: stocks.ExchangeParams{
				ConvertFrom: exchanger.USD,
				ConvertTo:   exchanger.UAH,
			},
			repoRes: stocks.PortfolioState{
				All: []stocks.PortfolioTickerState{
					{
						TickerID: "1",
						Name:     "aapl",
						Amount:   3,
						High:     2,
						Low:      2,
						Open:     2,
						Close:    2,
					},
				},
				Total: stocks.PortfolioTickerState{
					Amount: 3,
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
			expRes: stocks.PortfolioState{
				All: []stocks.PortfolioTickerState{
					{
						TickerID: "1",
						Name:     "aapl",
						Amount:   3,
						High:     4,
						Low:      4,
						Open:     4,
						Close:    4,
					},
				},
				Total: stocks.PortfolioTickerState{
					Amount: 3,
					High:   4,
					Low:    4,
					Open:   4,
					Close:  4,
				},
			},
		},
		"empty": {
			input: stocks.ExchangeParams{
				ConvertFrom: exchanger.USD,
				ConvertTo:   exchanger.UAH,
			},
			repoRes: stocks.PortfolioState{
				Total: stocks.PortfolioTickerState{},
			},
			expRes: stocks.PortfolioState{
				Total: stocks.PortfolioTickerState{},
			},
		},
	} {
		tCase := tCase
		t.Run(tName, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := mock.NewMockPortfolioRepository(ctrl)
			db := dbmock.NewMockDB(ctrl)
			exch := exchmock.NewMockAllCurrencyRater(ctrl)
			tr := mock.NewMockTickerRepository(ctrl)

			r.EXPECT().
				CountPortfolio(gomock.Any(), gomock.Any(), "1").
				Return(tCase.repoRes, nil)

			if tCase.exchangeRate != nil {
				exch.EXPECT().
					FindRates(gomock.Any(), gomock.Any()).
					Do(func(_ context.Context, params exchanger.ConvertCurrencyParams) {
						require.Equal(t, params.ConvertTo, tCase.input.ConvertTo)
						require.Equal(t, params.ConvertFrom, tCase.input.ConvertFrom)
					}).
					Return(tCase.exchangeRate, nil)
			}

			inter := NewPortfolioInteractor(db, r, exch, tr)
			res, err := inter.CountPortfolio(context.Background(), "1", tCase.input)

			require.NoError(t, err)
			require.Equal(t, tCase.expRes, res)
		})
	}
}
