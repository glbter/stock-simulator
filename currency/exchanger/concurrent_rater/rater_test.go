package rater

import (
	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/glbter/currency-ex/currency/exchanger/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRater_FindRates(t *testing.T) {
	now := time.Now()

	for tName, tData := range map[string]struct {
		expRes []exchanger.CurrencyRate
		expErr error
		errors []error
	}{
		"ok": {
			expRes: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now,
					Sale:     1,
					Purchase: 1,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now.Add(time.Hour),
					Sale:     1,
					Purchase: 1,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now.Add(2 * time.Hour),
					Sale:     1,
					Purchase: 1,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now.Add(3 * time.Hour),
					Sale:     1,
					Purchase: 1,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now.Add(4 * time.Hour),
					Sale:     1,
					Purchase: 1,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now.Add(5 * time.Hour),
					Sale:     1,
					Purchase: 1,
				},
			},
			errors: []error{
				nil, nil, nil, nil, nil, nil,
			},
		},
		"error": {
			expRes: []exchanger.CurrencyRate{
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now,
					Sale:     1,
					Purchase: 1,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now.Add(time.Hour),
					Sale:     1,
					Purchase: 1,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now.Add(2 * time.Hour),
					Sale:     1,
					Purchase: 1,
				},
				// removed by error
				//{
				//	Base: exchanger.UAH,
				//	Rated: exchanger.USD,
				//	Date: now.Add(3*time.Hour),
				//	Sale: 1,
				//	Purchase: 1,
				//},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now.Add(4 * time.Hour),
					Sale:     1,
					Purchase: 1,
				},
				{
					Base:     exchanger.UAH,
					Rated:    exchanger.USD,
					Date:     now.Add(5 * time.Hour),
					Sale:     1,
					Purchase: 1,
				},
			},
			errors: []error{
				nil, nil, nil, assert.AnError, nil, nil,
			},
			expErr: assert.AnError,
		},
	} {
		tData := tData
		t.Run(tName, func(t *testing.T) {
			//ctrl := gomock.NewController(t)
			//defer ctrl.Finish()
			//mr := mock.NewMockCurrencyRater(ctrl)
			//r := Rater{c: mr}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock.NewMockCurrencyRater(ctrl)
			r := Rater{
				c:            mr,
				timeInterval: time.Hour,
			}

			for i := 0; i < 6; i++ {
				t := now.Add(time.Duration(i) * time.Hour)
				mr.EXPECT().
					FindRate(exchanger.USD, t).
					Return(
						exchanger.CurrencyRate{
							Base:     exchanger.UAH,
							Rated:    exchanger.USD,
							Date:     t,
							Sale:     1,
							Purchase: 1,
						},
						tData.errors[i],
					)
			}

			res, err := r.FindRates(exchanger.USD, now, now.Add(6*time.Hour))

			require.Equal(t, tData.expRes, res)
			if tData.expErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tData.expErr.Error())
			}
		})
	}
}
