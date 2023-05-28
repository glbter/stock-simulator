package factory

import (
	"github.com/glbter/currency-ex/currency/exchanger"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
)

func NewDBCurrencyRater(
	db sqlc.DB,
) (exchanger.AllCurrencyRater, error) {
	rater, err := SeriesExchangerFactory(ExchangerFactoryParams{Db: db})
	if err != nil {
		return exchanger.AllCurrencyRater{}, err
	}

	return exchanger.NewAllCurrencyRater(rater), nil
}
