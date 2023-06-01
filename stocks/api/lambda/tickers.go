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

	tickerDaily, err := h.usecases.QueryLatestDaily(ctx,
		stocks.QueryDailyFilter{
			TickerIDs: tickerIDs,
		},
		stocks.ExchangeParams{
			ConverFrom: exchanger.USD,
			ConvertTo:  exchanger.USD,
		},
	)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, err
	}

	b, err := json.Marshal(tickerDaily)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
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
			}, fmt.Errorf("parse date before %w", err)
		}

		params.BeforeDateInc = &before
	}

	afterS, ok := ExtractFromQuery(request, afterInQuery)
	if ok {
		after, err := time.Parse(time.DateOnly, afterS)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusBadRequest,
			}, fmt.Errorf("parse date after %w", err)
		}

		params.AfterDateInc = &after
	}

	graph, err := h.usecases.QueryTickerDailyGraph(ctx, params, stocks.ExchangeParams{
		ConverFrom: exchanger.USD,
		ConvertTo:  exchanger.USD,
	})
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, err
	}

	b, err := json.Marshal(graph)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		Body:       string(b),
		StatusCode: http.StatusOK,
	}, nil
}
