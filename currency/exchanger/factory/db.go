package factory

import (
	sqlc "github.com/glbter/currency-ex/pkg/sql"
)

func NewDBCurrencyRater(
	db sqlc.DB,
) (AllCurrencyRater, error) {
	rater, err := SeriesExchangerFactory(ExchangerFactoryParams{Db: db})
	if err != nil {
		return AllCurrencyRater{}, err
	}

	return NewAllCurrencyRater(rater), nil
}
