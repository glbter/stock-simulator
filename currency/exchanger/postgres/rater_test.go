package postgres

import (
	"context"
	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRater_FindRates(t *testing.T) {
	for tName, tCase := range map[string]struct {
		params   FindExchangeRateParams
		expected []exchanger.CurrencyRate
	}{
		"source": {
			params: FindExchangeRateParams{
				Currency: exchanger.USD,
				Source:   exchanger.ExchangeSourcePrivat,
			},
			expected: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.1,
					Purchase: 8.2,
					Date:     time.Date(2010, 01, 22, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourcePrivat,
				},
			},
		},
		"full_interval": {
			params: FindExchangeRateParams{
				Currency: exchanger.USD,
				Source:   exchanger.ExchangeSourceNBU,
			},
			expected: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8,
					Purchase: 8,
					Date:     time.Date(2010, 01, 20, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.1,
					Purchase: 8.1,
					Date:     time.Date(2010, 01, 21, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.1,
					Purchase: 8.2,
					Date:     time.Date(2010, 01, 22, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.3,
					Purchase: 8.3,
					Date:     time.Date(2010, 01, 23, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.2,
					Purchase: 8.2,
					Date:     time.Date(2010, 01, 24, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
			},
		},
		"from_interval": {
			params: FindExchangeRateParams{
				Currency: exchanger.USD,
				Source:   exchanger.ExchangeSourceNBU,
				Start:    pointerTime(time.Date(2010, 01, 21, 0, 0, 0, 0, time.UTC)), //2010-01-24T00:00:00Z
			},
			expected: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.1,
					Purchase: 8.1,
					Date:     time.Date(2010, 01, 21, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.1,
					Purchase: 8.2,
					Date:     time.Date(2010, 01, 22, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.3,
					Purchase: 8.3,
					Date:     time.Date(2010, 01, 23, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.2,
					Purchase: 8.2,
					Date:     time.Date(2010, 01, 24, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
			},
		},
		"to_interval": {
			params: FindExchangeRateParams{
				Currency: exchanger.USD,
				Source:   exchanger.ExchangeSourceNBU,
				End:      pointerTime(time.Date(2010, 01, 23, 0, 0, 0, 0, time.UTC)), //2010-01-24T00:00:00Z
			},
			expected: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8,
					Purchase: 8,
					Date:     time.Date(2010, 01, 20, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.1,
					Purchase: 8.1,
					Date:     time.Date(2010, 01, 21, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.1,
					Purchase: 8.2,
					Date:     time.Date(2010, 01, 22, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.3,
					Purchase: 8.3,
					Date:     time.Date(2010, 01, 23, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
			},
		},
		"interval": {
			params: FindExchangeRateParams{
				Currency: exchanger.USD,
				Source:   exchanger.ExchangeSourceNBU,
				Start:    pointerTime(time.Date(2010, 01, 21, 0, 0, 0, 0, time.UTC)), //2010-01-24T00:00:00Z
				End:      pointerTime(time.Date(2010, 01, 23, 0, 0, 0, 0, time.UTC)), //2010-01-24T00:00:00Z
			},
			expected: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.1,
					Purchase: 8.1,
					Date:     time.Date(2010, 01, 21, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.1,
					Purchase: 8.2,
					Date:     time.Date(2010, 01, 22, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     8.3,
					Purchase: 8.3,
					Date:     time.Date(2010, 01, 23, 0, 0, 0, 0, time.UTC),
					Source:   exchanger.ExchangeSourceNBU,
				}},
		},
	} {
		tCase := tCase
		t.Run(tName, func(t *testing.T) {
			defer mustTruncateTables(t, testDB)

			fx.MustExecSQLFixture(t, testDB, "data.sql")

			r, err := Rater{}.
				FindRates(context.Background(), testDB, tCase.params)

			require.NoError(t, err)
			require.ElementsMatch(t, tCase.expected, r)
		})
	}
}
