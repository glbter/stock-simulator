package privat_bank

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

func TestPrivatBankCurrencyRater_FindRate(t *testing.T) {
	cTime := time.Now()

	for tName, tData := range map[string]struct {
		ReturnedData string
		ReturnedCode int
		ExpectedVal  exchanger.CurrencyRate
		ExpectedErr  error
	}{
		"success": {
			ReturnedData: `{
				"date": "01.12.2014",
				"bank": "PB",
				"baseCurrency": 980,
				"baseCurrencyLit": "UAH",
				"exchangeRate": [
					{
						"baseCurrency": "UAH",
						"currency": "USD",
						"saleRateNB": 15.0564130,
						"purchaseRateNB": 15.0564130,
						"saleRate": 15.7000000,
						"purchaseRate": 15.3500000
					},
					{
						"baseCurrency": "UAH",
						"currency": "EUR",
						"saleRateNB": 18.7949200,
						"purchaseRateNB": 18.7949200,
						"saleRate": 20.0000000,
						"purchaseRate": 19.2000000
					}
				]
			}`,
			ReturnedCode: http.StatusOK,
			ExpectedVal: exchanger.CurrencyRate{
				Base:     exchanger.UAH,
				Rated:    exchanger.USD,
				Sale:     15.7,
				Purchase: 15.35,
				Date:     cTime,
			},
		},
		"cant_find_in_response": {
			ReturnedData: `{
				"date": "01.12.2014",
				"bank": "PB",
				"baseCurrency": 980,
				"baseCurrencyLit": "UAH",
				"exchangeRate": [
					{
						"baseCurrency": "UAH",
						"currency": "EUR",
						"saleRateNB": 18.7949200,
						"purchaseRateNB": 18.7949200,
						"saleRate": 20.0000000,
						"purchaseRate": 19.2000000
					}
				]
			}`,
			ReturnedCode: http.StatusOK,
			ExpectedErr:  errors.New("not found"),
		},
		"success_empty_data_fields": {
			ReturnedData: `{
				"exchangeRate": [
					{
						"baseCurrency": "UAH",
						"currency": "USD",
						"saleRateNB": 15.0564130,
						"purchaseRateNB": 15.056413
					}
				]
			}`,
			ReturnedCode: http.StatusOK,
			ExpectedVal: exchanger.CurrencyRate{
				Base:     exchanger.UAH,
				Rated:    exchanger.USD,
				Sale:     15.056413,
				Purchase: 15.056413,
				Date:     cTime,
			},
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
