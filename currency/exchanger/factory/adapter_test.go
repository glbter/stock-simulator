package factory

import (
	"context"
	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/glbter/currency-ex/currency/exchanger/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAllCurrencyRater_FindRates(t *testing.T) {
	start := time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC)

	for tName, tCase := range map[string]struct {
		params    exchanger.ConvertCurrencyParams
		res       []exchanger.CurrencyRate
		expRes    []exchanger.CurrencyRate
		paramTo   exchanger.Currency
		err       error
		notMocked bool
	}{
		"convert_to_uah": {
			params: exchanger.ConvertCurrencyParams{
				ConvertFrom: exchanger.USD,
				ConvertTo:   exchanger.UAH,
				Start:       start,
				End:         start.Add(time.Hour * 24 * 7),
			},
			paramTo: exchanger.USD,
			res: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     2,
					Purchase: 2,
					Date:     start,
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     4,
					Purchase: 4,
					Date:     start.Add(time.Hour * 24),
					Source:   exchanger.ExchangeSourceNBU,
				},
			},
			expRes: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     2,
					Purchase: 2,
					Date:     start,
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     4,
					Purchase: 4,
					Date:     start.Add(time.Hour * 24),
					Source:   exchanger.ExchangeSourceNBU,
				},
			},
		},
		"convert_from_uah": {
			params: exchanger.ConvertCurrencyParams{
				ConvertFrom: exchanger.UAH,
				ConvertTo:   exchanger.USD,
				Start:       start,
				End:         start.Add(time.Hour * 24 * 7),
			},
			paramTo: exchanger.USD,
			res: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     2,
					Purchase: 2,
					Date:     start,
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     4,
					Purchase: 4,
					Date:     start.Add(time.Hour * 24),
					Source:   exchanger.ExchangeSourceNBU,
				},
			},
			expRes: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     0.5,
					Purchase: 0.5,
					Date:     start,
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Sale:     0.25,
					Purchase: 0.25,
					Date:     start.Add(time.Hour * 24),
					Source:   exchanger.ExchangeSourceNBU,
				},
			},
		},
		"equal_currency": {
			notMocked: true,
			params: exchanger.ConvertCurrencyParams{
				ConvertFrom: exchanger.USD,
				ConvertTo:   exchanger.USD,
				Start:       start,
				End:         start.Add(time.Hour * 24),
			},
			paramTo: exchanger.USD,
			expRes: []exchanger.CurrencyRate{
				{
					Base:     exchanger.USD,
					Rated:    exchanger.USD,
					Sale:     1,
					Purchase: 1,
					Date:     start,
					Source:   exchanger.ExchangeSourceNBU,
				},
				{
					Base:     exchanger.USD,
					Rated:    exchanger.USD,
					Sale:     1,
					Purchase: 1,
					Date:     start.Add(time.Hour * 24),
					Source:   exchanger.ExchangeSourceNBU,
				},
			},
		},
	} {
		tCase := tCase
		t.Run(tName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cr := mock.NewMockCurrencySeriesRater(ctrl)
			if !tCase.notMocked {
				cr.EXPECT().
					FindRates(gomock.Any(), tCase.paramTo, tCase.params.Start, tCase.params.End).
					Return(tCase.res, tCase.err)
			}

			acr := NewAllCurrencyRater(cr)

			rates, err := acr.FindRates(context.Background(), tCase.params)
			require.NoError(t, err)
			require.ElementsMatch(t, rates, tCase.expRes)
		})
	}
}
