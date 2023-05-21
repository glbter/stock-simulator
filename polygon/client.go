package polygon

import (
	"context"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"net/http"
)

func New(apiKey string, client *http.Client) {
	pc := polygon.NewWithClient(apiKey, client)

	ctx := context.Background()

	pc.ListTickers(ctx, &models.ListTickersParams{})

	pc.GetMarketHolidays(ctx)

	pc.ListSplits(ctx, &models.ListSplitsParams{})
}
