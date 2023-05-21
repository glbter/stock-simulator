package nbu

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/stretchr/testify/require"
)

func TestNBUCurrencyRater_FindRate(t *testing.T) {
	cTime := time.Now()

	for tName, tData := range map[string]struct {
		ReturnedData string
		ReturnedCode int
		ExpectedVal  exchanger.CurrencyRate
		ExpectedErr  error
	}{
		"success": {
			ReturnedData: `[
				{
					"r030": 826,
					"txt": "Фунт стерлінгів",
					"rate": 11.918468,
					"cc": "GBP",
					"exchangedate": "02.03.2010"
				},
				{
					"r030": 840,
					"txt": "Долар США",
					"rate": 7.99,
					"cc": "USD",
					"exchangedate": "02.03.2010"
				}
					]`,
			ReturnedCode: http.StatusOK,
			ExpectedVal: exchanger.CurrencyRate{
				Base:     exchanger.UAH,
				Rated:    exchanger.USD,
				Sale:     7.99,
				Purchase: 7.99,
				Date:     cTime,
			},
		},
		"cant_find_in_response": {
			ReturnedData: `[
				{
					"r030": 826,
					"txt": "Фунт стерлінгів",
					"rate": 11.918468,
					"cc": "GBP",
					"exchangedate": "02.03.2010"
				}
					]`,
			ReturnedCode: http.StatusOK,
			ExpectedErr:  errors.New("not found"),
		},
		"corrupt_response": {
			ReturnedData: `{}`,
			ReturnedCode: http.StatusOK,
			ExpectedErr:  errors.New("json: cannot unmarshal object"),
		},
		"not_found_in_response": {
			ReturnedData: `{}`,
			ReturnedCode: http.StatusNotFound,
			ExpectedErr:  errors.New("responded with 404 http code"),
		},
	} {
		tData := tData
		t.Run(tName, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tData.ReturnedCode)
				fmt.Fprintf(w, tData.ReturnedData)
			}))
			defer svr.Close()

			resp, err := NewClient(http.DefaultClient, svr.URL).
				FindRate(exchanger.USD, cTime)

			if tData.ExpectedErr != nil {
				require.Contains(t, err.Error(), tData.ExpectedErr.Error())
			}

			require.Equal(t, tData.ExpectedVal, resp)
		})
	}
}
