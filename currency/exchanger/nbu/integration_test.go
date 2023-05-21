//go:integration_test

package nbu

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/stretchr/testify/require"
)

func TestNBUCurrencyRater_FindRate_Integration(t *testing.T) {
	url := os.Getenv("NBU_CURRENCY_URL")
	if url == "" {
		t.Skip("set NBU_CURRENCY_URL to run integration test")
	}

	cTime := time.Date(2010, time.March, 02, 0, 0, 0, 0, time.UTC)

	resp, err := NewClient(
		&http.Client{
			Timeout: time.Second,
		},
		url,
	).
		FindRate(exchanger.USD, cTime)

	require.NoError(t, err)
	require.Equal(t,
		exchanger.CurrencyRate{
			Base:     exchanger.UAH,
			Rated:    exchanger.USD,
			Sale:     8.0004,
			Purchase: 8.0004,
			Date:     cTime,
		},
		resp,
	)
}
