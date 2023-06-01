package usecases

import (
	"context"
	"github.com/glbter/currency-ex/currency/exchanger"
	exchmock "github.com/glbter/currency-ex/currency/exchanger/mock"
	dbmock "github.com/glbter/currency-ex/pkg/sql/mock"
	"github.com/glbter/currency-ex/stocks"
	"github.com/glbter/currency-ex/stocks/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestPortfolioInteractor_TradeTickers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock.NewMockPortfolioRepository(ctrl)
	db := dbmock.NewMockDB(ctrl)
	exch := exchmock.NewMockAllCurrencyRater(ctrl)
	tr := mock.NewMockTickerRepository(ctrl)

	inter := NewPortfolioInteractor(db, r, exch, tr)
	//inter := NewTickerInteractor(db, r)
	inter.CountPortfolio(context.Background(), "1", stocks.ExchangeParams{
		ConverFrom: exchanger.USD,
		ConvertTo:  exchanger.UAH,
	})

}
