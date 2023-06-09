package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/glbter/currency-ex/pkg/serrors"
	"github.com/glbter/currency-ex/stocks"
	"net/http"
	"time"
)

const (
	tickerNamesInQuery = "ticker_names"
)

type TickerHandler struct {
	usecases stocks.TickerUsecases
}

func NewTickerHandler(
	usecases stocks.TickerUsecases,
) TickerHandler {
	return TickerHandler{
		usecases: usecases,
	}
}

func (h TickerHandler) GetTickers(
	ctx context.Context,
	request events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	tickerIDs, ok := ExtractSliceFromQuery(request, tickerIDsInQuery)
	if !ok {
		tickerIDs = nil
	}

	tickerNames, ok := ExtractSliceFromQuery(request, tickerNamesInQuery)
	if !ok {
		tickerNames = nil
	}

	tickerDaily, err := h.usecases.QueryLatestDaily(ctx,
		stocks.QueryDailyFilter{
			TickerIDs: tickerIDs,
			Tickers:   tickerNames,
		},
		stocks.ExchangeParams{
			ConvertFrom: exchanger.USD,
			ConvertTo:   exchanger.UAH,
		},
	)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, serrors.GetErrorByTypeAndLog(err)
	}

	b, err := json.Marshal(tickerDaily)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, serrors.GetErrorByTypeAndLog(err)
	}

	return events.APIGatewayV2HTTPResponse{
		Body:       string(b),
		StatusCode: http.StatusOK,
	}, nil
}

func (h TickerHandler) GetTickerGraph(
	ctx context.Context,
	request events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	tickerIDs, ok := ExtractSliceFromQuery(request, tickerIDsInQuery)
	if !ok {
		tickerIDs = nil
	}

	params := stocks.QueryDailyGraphParams{
		TickerIDs: tickerIDs,
	}

	beforeS, ok := ExtractFromQuery(request, beforeInQuery)
	if ok {
		before, err := time.Parse(time.DateOnly, beforeS)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusBadRequest,
			}, serrors.GetErrorByTypeAndLog(fmt.Errorf("%w: parse date before: %v", serrors.ErrBadInput, err))
		}

		params.BeforeDateInc = &before
	}

	afterS, ok := ExtractFromQuery(request, afterInQuery)
	if ok {
		after, err := time.Parse(time.DateOnly, afterS)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusBadRequest,
			}, serrors.GetErrorByTypeAndLog(fmt.Errorf("%w: parse date after: %v", serrors.ErrBadInput, err))
		}

		params.AfterDateInc = &after
	}

	graph, err := h.usecases.QueryTickerDailyGraph(ctx, params, stocks.ExchangeParams{
		ConvertFrom: exchanger.USD,
		ConvertTo:   exchanger.USD,
	})
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, serrors.GetErrorByTypeAndLog(err)
	}

	b, err := json.Marshal(graph)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, serrors.GetErrorByTypeAndLog(err)
	}

	return events.APIGatewayV2HTTPResponse{
		Body:       string(b),
		StatusCode: http.StatusOK,
	}, nil
}
