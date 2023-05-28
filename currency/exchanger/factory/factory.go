package factory

import (
	"fmt"
	"github.com/glbter/currency-ex/currency/exchanger"
	concurent "github.com/glbter/currency-ex/currency/exchanger/concurrent"
	"github.com/glbter/currency-ex/currency/exchanger/nbu"
	"github.com/glbter/currency-ex/currency/exchanger/privat_bank"
	"github.com/glbter/currency-ex/pkg/serrors"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	"net/http"
	"time"

	"github.com/glbter/currency-ex/currency/exchanger/postgres"
)

type ExchangerFactoryParams struct {
	Db     sqlc.DB
	Source exchanger.ExchangeSource
	Url    string
}

func SeriesExchangerFactory(p ExchangerFactoryParams) (exchanger.CurrencySeriesRater, error) {
	if p.Db != nil {
		return postgres.NewCurryingExchangeAdapter(p.Db, postgres.Rater{}, p.Source), nil
	}

	c := &http.Client{Timeout: 2 * time.Second}
	if p.Source == exchanger.ExchangeSourcePrivat {
		return concurent.NewRater(privat_bank.NewClient(c, p.Url), time.Hour*24), nil
	}

	if p.Source == exchanger.ExchangeSourceNBU {
		return concurent.NewRater(nbu.NewClient(c, p.Url), time.Hour*24), nil
	}

	return nil, fmt.Errorf("%w: %v", serrors.ErrBadInput, "no such provider")
}
